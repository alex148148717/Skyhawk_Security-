package domain

type AverageTeamSeason struct {
	SeasonID         int
	TeamID           int
	TeamName         string
	NbaTeamID        int
	AvgPoints        float64
	AvgRebounds      float64
	AvgAssists       float64
	AvgSteals        float64
	AvgBlocks        float64
	AvgFouls         float64
	AvgTurnovers     float64
	AvgMinutesPlayed float64
}
