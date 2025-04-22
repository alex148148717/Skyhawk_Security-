package team_player_season_statistics

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
)

type DBClient struct {
	db *sql.DB
}

func NewDBClient(db *sql.DB) *DBClient {
	c := &DBClient{db: db}
	return c
}

func (c *DBClient) AveragePlayerSeason(ctx context.Context, ids []int32) ([]AveragePlayerSeason, error) {

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
	var results []AveragePlayerSeason
	for rows.Next() {
		var r AveragePlayerSeason
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
