package infrastructure

import "fmt"

func KeyGenerate(seasonID, playerID int) string {
	return fmt.Sprintf("AveragePlayerSeasonClient_%d_%d", seasonID, playerID)

}
