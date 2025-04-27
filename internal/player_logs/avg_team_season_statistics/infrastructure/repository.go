package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/lib/pq"
	"skyhawk/internal/player_logs/avg_team_season_statistics/domain"
	"skyhawk/internal/player_logs/config"
)

type Repository struct {
	db              *sql.DB
	dynamoDBAPI     dynamodbiface.DynamoDBAPI
	dynamoTableName string
}

func NewRepository(db *sql.DB, config *config.Config, dynamoDBAPI dynamodbiface.DynamoDBAPI) domain.TeamGameStatisticRepository {
	r := Repository{
		db:              db,
		dynamoDBAPI:     dynamoDBAPI,
		dynamoTableName: config.CacheTableName,
	}
	return &r
}

func (c *Repository) GetAverageTeamsSeason(ctx context.Context, ids []int32) ([]domain.AverageTeamSeason, error) {

	query := `
WITH grouped_stats AS (
  SELECT
    season_id,
    team_id,
    player_id,
    ROUND(SUM(points)::numeric / COUNT(DISTINCT game_id), 2) AS avg_points,
    ROUND(SUM(rebounds)::numeric / COUNT(DISTINCT game_id), 2) AS avg_rebounds,
    ROUND(SUM(assists)::numeric / COUNT(DISTINCT game_id), 2) AS avg_assists,
    ROUND(SUM(steals)::numeric / COUNT(DISTINCT game_id), 2) AS avg_steals,
    ROUND(SUM(blocks)::numeric / COUNT(DISTINCT game_id), 2) AS avg_blocks,
    ROUND(SUM(fouls)::numeric / COUNT(DISTINCT game_id), 2) AS avg_fouls,
    ROUND(SUM(turnovers)::numeric / COUNT(DISTINCT game_id), 2) AS avg_turnovers,
    SUM(minutes_played)::numeric / COUNT(DISTINCT game_id) AS avg_minutes_played
  FROM player_stats_raw
  WHERE (season_id, team_id) IN (
    SELECT DISTINCT season_id, team_id
    FROM player_stats_raw
    WHERE id = ANY($1)
  )

  GROUP BY season_id, team_id, player_id
), team_avg as (

SELECT gs.season_id,
                       gs.team_id,

                       sum(gs.avg_points)/ COUNT(DISTINCT player_id)as avg_points,
                         sum(gs.avg_rebounds)/ COUNT(DISTINCT player_id)as avg_rebounds,
                         sum(gs.avg_assists)/ COUNT(DISTINCT player_id)as avg_assists,
                         sum(gs.avg_steals)/ COUNT(DISTINCT player_id)as avg_steals,
                         sum(gs.avg_blocks)/ COUNT(DISTINCT player_id)as avg_blocks,
                         sum(gs.avg_fouls)/ COUNT(DISTINCT player_id)as avg_fouls,
                         sum(gs.avg_turnovers)/ COUNT(DISTINCT player_id)as avg_turnovers,
                         sum(gs.avg_minutes_played)/ COUNT(DISTINCT player_id)as avg_minutes_played
                FROM grouped_stats gs
          group by season_id,team_id

                )

SELECT
  team_avg.season_id,
  team_avg.team_id,
  t.name AS team_name,
  t.nba_team_id,
  team_avg.avg_points,
  team_avg.avg_rebounds,
  team_avg.avg_assists,
  team_avg.avg_steals,
  team_avg.avg_blocks,
  team_avg.avg_fouls,
  team_avg.avg_turnovers,
  team_avg.avg_minutes_played
FROM team_avg
JOIN teams t ON team_id = t.id

where season_id=2024 and team_avg.team_id=1
`
	rows, err := c.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []domain.AverageTeamSeason
	for rows.Next() {
		var r domain.AverageTeamSeason
		if err := rows.Scan(
			&r.SeasonID,
			&r.TeamID,
			&r.TeamName,
			&r.NbaTeamID,
			&r.AvgPoints,
			&r.AvgRebounds,
			&r.AvgAssists,
			&r.AvgSteals,
			&r.AvgBlocks,
			&r.AvgFouls,
			&r.AvgTurnovers,
			&r.AvgMinutesPlayed,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}
func (c *Repository) GetCacheAverageTeamSeason(ctx context.Context, seasonID, teamID int) (*domain.AverageTeamSeason, error) {

	key := KeyGenerate(seasonID, teamID)
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.dynamoTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
	}

	result, err := c.dynamoDBAPI.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("key not found")
	}
	var out domain.AverageTeamSeason
	err = dynamodbattribute.UnmarshalMap(result.Item, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal struct: %w", err)
	}
	return &out, nil
}
func (c *Repository) SetCacheAverageTeamSeason(ctx context.Context, seasonID, teamID int, averageTeamSeason domain.AverageTeamSeason) error {
	item, err := dynamodbattribute.MarshalMap(averageTeamSeason)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}
	key := KeyGenerate(seasonID, teamID)

	item["id"] = &dynamodb.AttributeValue{
		S: aws.String(key),
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.dynamoTableName),
		Item:      item,
	}

	_, err = c.dynamoDBAPI.PutItem(input)
	return err
}
