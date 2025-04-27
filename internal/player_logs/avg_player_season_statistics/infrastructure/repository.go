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
	"skyhawk/internal/player_logs/avg_player_season_statistics/domain"
	"skyhawk/internal/player_logs/config"
)

type Repository struct {
	db              *sql.DB
	dynamoDBAPI     dynamodbiface.DynamoDBAPI
	dynamoTableName string
}

func NewRepository(db *sql.DB, config *config.Config, dynamoDBAPI dynamodbiface.DynamoDBAPI) domain.PlayerGameStatisticRepository {
	r := Repository{
		db:              db,
		dynamoDBAPI:     dynamoDBAPI,
		dynamoTableName: config.CacheTableName,
	}
	return &r
}

func (c *Repository) GetAveragePlayersSeason(ctx context.Context, ids []int32) ([]domain.AveragePlayerSeason, error) {

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
  WHERE (season_id, team_id, player_id) IN (
    SELECT DISTINCT season_id, team_id, player_id
    FROM player_stats_raw
    WHERE id = ANY($1)
  )
  GROUP BY season_id, team_id, player_id
)
SELECT
  gs.season_id,
  gs.team_id,
  t.name AS team_name,
  t.nba_team_id AS nba_team_id,
  gs.player_id,
  p.name AS player_name,
  p.nba_player_id AS nba_player_id,
  ptn.jersey_number,
  gs.avg_points,
  gs.avg_rebounds,
  gs.avg_assists,
  gs.avg_steals,
  gs.avg_blocks,
  gs.avg_fouls,
  gs.avg_turnovers,
  gs.avg_minutes_played
FROM grouped_stats gs
JOIN player_team_numbers ptn
  ON gs.player_id = ptn.player_id
 AND gs.team_id = ptn.team_id
 AND gs.season_id = ptn.season_year
JOIN players p ON gs.player_id = p.id
JOIN teams t ON gs.team_id = t.id
`
	rows, err := c.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []domain.AveragePlayerSeason
	for rows.Next() {
		var r domain.AveragePlayerSeason
		if err := rows.Scan(
			&r.SeasonID,
			&r.TeamID,
			&r.TeamName,
			&r.NbaTeamID,
			&r.PlayerID,
			&r.PlayerName,
			&r.NbaPlayerID,
			&r.JerseyNumber,
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
func (c *Repository) GetCacheAveragePlayerSeason(ctx context.Context, seasonID, playerID int) (*domain.AveragePlayerSeason, error) {

	key := KeyGenerate(seasonID, playerID)
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
	var out domain.AveragePlayerSeason
	err = dynamodbattribute.UnmarshalMap(result.Item, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal struct: %w", err)
	}
	return &out, nil
}
func (c *Repository) SetCacheAveragePlayerSeason(ctx context.Context, seasonID, playerID int, averagePlayerSeason domain.AveragePlayerSeason) error {
	item, err := dynamodbattribute.MarshalMap(averagePlayerSeason)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}
	key := KeyGenerate(seasonID, playerID)

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
