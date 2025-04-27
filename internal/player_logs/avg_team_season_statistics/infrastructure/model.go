package infrastructure

import "fmt"

func KeyGenerate(seasonID, teamID int) string {
	return fmt.Sprintf("AverageTeamSeasonClient_%d_%d", seasonID, teamID)

}
