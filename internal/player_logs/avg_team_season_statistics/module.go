package avg_team_season_statistics

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"skyhawk/internal/player_logs/avg_team_season_statistics/domain"
	"skyhawk/internal/player_logs/avg_team_season_statistics/infrastructure"
	"skyhawk/internal/player_logs/avg_team_season_statistics/interfaces"
)

var Module = fx.Options(
	fx.Provide(
		domain.NewService,
		infrastructure.NewRepository,
		interfaces.NewHandler,
		interfaces.NewPageTemplate,
	),
	fx.Invoke(func(router *chi.Mux, handler *interfaces.AVGTeamGameStatisticHandler) {
		router.Get("/log_game_player_statistic/v1/season/{SeasonID}/team/{NbaTeamID}", handler.Handler)

	}),
)
