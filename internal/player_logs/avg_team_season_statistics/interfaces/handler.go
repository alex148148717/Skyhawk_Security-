package interfaces

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"skyhawk/internal/player_logs/avg_team_season_statistics/domain"
	"strconv"
)

type AVGTeamGameStatisticHandler struct {
	teamGameStatisticService domain.Service
	pageTemplate             *PageTemplate
}

func NewHandler(teamGameStatisticService domain.Service, pageTemplate *PageTemplate) *AVGTeamGameStatisticHandler {
	h := AVGTeamGameStatisticHandler{teamGameStatisticService: teamGameStatisticService, pageTemplate: pageTemplate}
	return &h
}

func (c *AVGTeamGameStatisticHandler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := strconv.Atoi(chi.URLParam(r, "SeasonID"))
	if err != nil {
		http.Error(w, "Invalid SeasonID", http.StatusBadRequest)
		return
	}

	nbaTeamID, err := strconv.Atoi(chi.URLParam(r, "NbaTeamID"))
	if err != nil {
		http.Error(w, "Invalid TeamID", http.StatusBadRequest)
		return
	}
	averageTeamSeason, err := c.teamGameStatisticService.GetTeamData(ctx, seasonID, nbaTeamID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.WriteHeader(http.StatusOK)
	if err := c.pageTemplate.pageTpl.Execute(w, averageTeamSeason); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
	return

}
