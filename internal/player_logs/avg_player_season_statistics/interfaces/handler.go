package interfaces

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"skyhawk/internal/player_logs/avg_player_season_statistics/domain"
	"strconv"
)

type AVGPlayerGameStatisticHandler struct {
	playerGameStatisticService domain.Service
	pageTemplate               *PageTemplate
}

func NewHandler(playerGameStatisticService domain.Service, pageTemplate *PageTemplate) *AVGPlayerGameStatisticHandler {
	h := AVGPlayerGameStatisticHandler{playerGameStatisticService: playerGameStatisticService, pageTemplate: pageTemplate}
	return &h
}

func (c *AVGPlayerGameStatisticHandler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := strconv.Atoi(chi.URLParam(r, "SeasonID"))
	if err != nil {
		http.Error(w, "Invalid SeasonID", http.StatusBadRequest)
		return
	}

	nbaPlayerID, err := strconv.Atoi(chi.URLParam(r, "NbaPlayerID"))
	if err != nil {
		http.Error(w, "Invalid TeamID", http.StatusBadRequest)
		return
	}
	averagePlayerSeason, err := c.playerGameStatisticService.GetPlayerData(ctx, seasonID, nbaPlayerID)
	_ = averagePlayerSeason

	if err != nil {
		fmt.Printf("failed to get cache item: %s\n", err)
		http.Error(w, "no data", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.WriteHeader(http.StatusOK)
	if err := c.pageTemplate.pageTpl.Execute(w, averagePlayerSeason); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
	return

}
