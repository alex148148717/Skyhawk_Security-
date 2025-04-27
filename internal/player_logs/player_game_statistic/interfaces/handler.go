package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"skyhawk/internal/player_logs/player_game_statistic/domain"
	"sync"
)

type PlayerGameStatisticHandler struct {
	playerGameStatisticService domain.PlayerGameStatisticService
}

func NewHandler(playerGameStatisticService domain.PlayerGameStatisticService) *PlayerGameStatisticHandler {
	h := PlayerGameStatisticHandler{playerGameStatisticService: playerGameStatisticService}
	return &h
}

func (c *PlayerGameStatisticHandler) RegisterRoutes(router chi.Router) {
	router.Post("/import/{ID}", c.Handler)

}

func (c *PlayerGameStatisticHandler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()
	var playerLogStatistics []PlayerLogStatisticRaw
	if err := json.NewDecoder(r.Body).Decode(&playerLogStatistics); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	JobID := chi.URLParam(r, "ID")

	workers := 100
	wg := &sync.WaitGroup{}
	mutex := sync.Mutex{}
	playerLogStatisticChan := make(chan PlayerLogStatisticRaw, workers)
	wg.Add(workers)
	var errors []error
	playerLogStatisticDB := make([]domain.PlayerGameStatistic, 0, len(playerLogStatistics))
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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ValidationErrorsResponse{Success: false, Errors: errors})
		return
	}
	err := c.playerGameStatisticService.InsertLines(ctx, JobID, playerLogStatisticDB)

	w.WriteHeader(http.StatusOK)

	if err != nil {
		json.NewEncoder(w).Encode(ImportResponse{Success: false, Errors: err})
		return
	}

	json.NewEncoder(w).Encode(ImportResponse{Success: true})

	return

}
