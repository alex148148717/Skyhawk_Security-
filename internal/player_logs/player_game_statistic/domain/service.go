package domain

import (
	"context"
	"skyhawk/internal/player_logs/avg_player_season_statistics/domain"
	avgteamDomain "skyhawk/internal/player_logs/avg_team_season_statistics/domain"
	"sync"
)

type PlayerGameStatisticRepository interface {
	InsertPlayerGameStatistic(ctx context.Context, JobID string, playerGameStatistics []PlayerGameStatistic) ([]int32, error)
}

type SyncClient interface {
	Sync(ctx context.Context, ids []int32) error
}

type SyncClients []SyncClient

func NewSyncClients(avgPlayerSeasonStatistics domain.Service, avgTeamSeasonStatistics avgteamDomain.Service) SyncClients {
	a := []SyncClient{avgPlayerSeasonStatistics, avgTeamSeasonStatistics}
	return a
}

type PlayerGameStatisticServiceImpl struct {
	playerGameStatisticRepository PlayerGameStatisticRepository
	syncClients                   SyncClients
}
type PlayerGameStatisticService interface {
	InsertLines(ctx context.Context, JobID string, lines []PlayerGameStatistic) error
}

func NewPlayerGameStatisticService(playerGameStatisticRepository PlayerGameStatisticRepository, syncClients SyncClients) PlayerGameStatisticService {
	p := PlayerGameStatisticServiceImpl{playerGameStatisticRepository: playerGameStatisticRepository, syncClients: syncClients}
	return &p
}
func (c *PlayerGameStatisticServiceImpl) InsertLines(ctx context.Context, JobID string, lines []PlayerGameStatistic) error {

	ids, err := c.playerGameStatisticRepository.InsertPlayerGameStatistic(ctx, JobID, lines)
	if err != nil {
		return err
	}

	syncClients := c.syncClients
	numWorkers := 4
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error, 1)
	clientCh := make(chan SyncClient)

	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for client := range clientCh {
				if err := client.Sync(ctx, ids); err != nil {
					select {
					case errCh <- err:
						cancel()
					default:
					}
					return
				}
			}
		}()
	}

	go func() {
		for _, client := range syncClients {
			select {
			case <-ctx.Done():
				break
			default:
				clientCh <- client
			}
		}
		close(clientCh)
	}()

	wg.Wait()
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
