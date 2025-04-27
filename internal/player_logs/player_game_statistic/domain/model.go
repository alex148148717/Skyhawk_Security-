package domain

type PlayerGameStatistic struct {
	ID            int32
	JobID         string
	SeasonID      int32
	GameID        int32
	TeamID        int32
	PlayerID      int32
	Points        int16
	Rebounds      int16
	Assists       int16
	Steals        int16
	Blocks        int16
	Fouls         uint8
	Turnovers     int16
	MinutesPlayed float32
}
