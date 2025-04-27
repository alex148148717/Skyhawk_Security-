package interfaces

import (
	"fmt"
	"skyhawk/internal/player_logs/player_game_statistic/domain"
)

var (
	FoulsError         = fmt.Errorf("invalid fouls (must be between 0–6)")
	MinutesPlayedError = fmt.Errorf("invalid minutes played (must be between 0–48.0)")
)

type PlayerLogStatisticRaw struct {
	Id            int32   `json:"id"`
	SeasonYear    int32   `json:"season_year"`
	GameId        int32   `json:"game_id"`
	TeamId        int32   `json:"team_id"`
	PlayerId      int32   `json:"player_id"`
	Points        int16   `json:"points"`
	Rebounds      int16   `json:"rebounds"`
	Assists       int16   `json:"assists"`
	Steals        int16   `json:"steals"`
	Blocks        int16   `json:"blocks"`
	Fouls         uint8   `json:"fouls"`
	Turnovers     int16   `json:"turnovers"`
	MinutesPlayed float32 `json:"minutes_played"`
}

func (p PlayerLogStatisticRaw) Convert() (*domain.PlayerGameStatistic, error) {

	if p.Fouls < 0 || p.Fouls > 6 {
		return nil, FoulsError
	}
	if p.MinutesPlayed < 0 || p.MinutesPlayed > 48.0 {
		return nil, MinutesPlayedError
	}

	playerGameStatistic := domain.PlayerGameStatistic{
		ID:            p.Id,
		SeasonID:      p.SeasonYear,
		GameID:        p.GameId,
		TeamID:        p.TeamId,
		PlayerID:      p.PlayerId,
		Points:        p.Points,
		Rebounds:      p.Rebounds,
		Assists:       p.Assists,
		Steals:        p.Steals,
		Blocks:        p.Blocks,
		Fouls:         p.Fouls,
		Turnovers:     p.Turnovers,
		MinutesPlayed: p.MinutesPlayed,
	}
	return &playerGameStatistic, nil
}

type ValidationErrorsResponse struct {
	Success bool    `json:"success"`
	Errors  []error `json:"errors"`
}

type ImportResponse struct {
	Success bool  `json:"success"`
	Errors  error `json:"error"`
}
