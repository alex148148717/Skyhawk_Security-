package infrastructure

import (
	"context"
	"database/sql"
	"skyhawk/internal/player_logs/player_game_statistic/domain"
)

type PlayerGameStatisticRepository struct {
	db *sql.DB
}

func NewPlayerGameStatisticRepository(db *sql.DB) domain.PlayerGameStatisticRepository {
	c := PlayerGameStatisticRepository{db: db}
	return &c
}

func (c *PlayerGameStatisticRepository) InsertPlayerGameStatistic(ctx context.Context, JobID string, playerGameStatistics []domain.PlayerGameStatistic) ([]int32, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`
	INSERT INTO player_stats_raw (
		id, job_id, season_id, game_id, team_id, player_id,
		points, rebounds, assists, steals, blocks, fouls,
		turnovers, minutes_played
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	ON CONFLICT (id) DO NOTHING
`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	ids := make(map[int32]bool, len(playerGameStatistics))
	for _, stat := range playerGameStatistics {
		_, err := stmt.Exec(
			stat.ID, JobID, stat.SeasonID, stat.GameID, stat.TeamID, stat.PlayerID,
			stat.Points, stat.Rebounds, stat.Assists, stat.Steals, stat.Blocks, stat.Fouls,
			stat.Turnovers, stat.MinutesPlayed,
		)
		if err != nil {
			return nil, err
		}
		ids[stat.ID] = true
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return MapValuesToSlice(ids), nil
}
func MapValuesToSlice[K comparable, V any](m map[K]V) []K {
	values := make([]K, 0, len(m))
	for k := range m {
		values = append(values, k)
	}
	return values
}
