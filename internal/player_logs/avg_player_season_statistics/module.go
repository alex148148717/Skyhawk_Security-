package avg_player_season_statistics

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"skyhawk/internal/player_logs/avg_player_season_statistics/domain"
	"skyhawk/internal/player_logs/avg_player_season_statistics/infrastructure"
	"skyhawk/internal/player_logs/avg_player_season_statistics/interfaces"
)

var Module = fx.Options(
	fx.Provide(
		domain.NewService,
		infrastructure.NewRepository,
		interfaces.NewHandler,
		interfaces.NewPageTemplate,
	),
	fx.Invoke(func(router *chi.Mux, handler *interfaces.AVGPlayerGameStatisticHandler) {
		router.Get("/log_game_player_statistic/v1/season/{SeasonID}/player/{NbaPlayerID}", handler.Handler)

	}),
)
