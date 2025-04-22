package player_game_statistic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"
)

type PlayerGameStatisticClient interface {
	InsertPlayerGameStatistic(ctx context.Context, JobID string, playerGameStatistics []PlayerGameStatistic) ([]int32, error)
}
type SyncClient interface {
	Sync(ctx context.Context, ids []int32) error
}

type Client struct {
	playerGameStatisticClient PlayerGameStatisticClient
	syncClients               []SyncClient
}

func New(r chi.Router, playerGameStatisticClient PlayerGameStatisticClient, syncClients []SyncClient) *Client {
	c := Client{playerGameStatisticClient: playerGameStatisticClient, syncClients: syncClients}
	r.Post("/import/{ID}", c.logGamePlayerStatisticsHandler)

	return &c
}

func (c *Client) logGamePlayerStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var playerLogStatistics []PlayerLogStatisticRaw
	if err := json.NewDecoder(r.Body).Decode(&playerLogStatistics); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	JobID := chi.URLParam(r, "ID")

	workers := 100
	wg := &sync.WaitGroup{}
	mutex := sync.Mutex{}
	playerLogStatisticChan := make(chan PlayerLogStatisticRaw, workers)
	wg.Add(workers)
	var errors []error
	playerLogStatisticDB := make([]PlayerGameStatistic, 0, len(playerLogStatistics))
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for playerLogStatistic := range playerLogStatisticChan {
				psdb, err := playerLogStatistic.Convert()
				mutex.Lock()
				if err != nil {
					errors = append(errors, fmt.Errorf("playerLogStatistic %v playerLogStatistic.Convert: %w", playerLogStatistic, err))
					mutex.Unlock()
					continue
				}
				playerLogStatisticDB = append(playerLogStatisticDB, *psdb)
				mutex.Unlock()
			}

		}()
	}
	for _, playerLogStatistic := range playerLogStatistics {
		playerLogStatisticChan <- playerLogStatistic
	}
	close(playerLogStatisticChan)
	wg.Wait()

	if len(errors) > 0 {
		json.NewEncoder(w).Encode(ValidationErrorsResponse{Success: false, Errors: errors})
		return
	}
	err := c.InsertLines(ctx, JobID, playerLogStatisticDB)
	if err != nil {

	}
	if err != nil {
		json.NewEncoder(w).Encode(ImportResponse{Success: false, Errors: err})
		return
	}
	json.NewEncoder(w).Encode(ImportResponse{Success: true})

	w.WriteHeader(http.StatusOK)
	return

}

func (c *Client) InsertLines(ctx context.Context, JobID string, lines []PlayerGameStatistic) error {

	ids, err := c.playerGameStatisticClient.InsertPlayerGameStatistic(ctx, JobID, lines)
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
