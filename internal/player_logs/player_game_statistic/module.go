package player_game_statistic

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"skyhawk/internal/player_logs/player_game_statistic/domain"
	"skyhawk/internal/player_logs/player_game_statistic/infrastructure"
	"skyhawk/internal/player_logs/player_game_statistic/interfaces"
)

var Module = fx.Options(
	fx.Provide(
		infrastructure.NewPlayerGameStatisticRepository,
		domain.NewPlayerGameStatisticService,
		domain.NewSyncClients,
		interfaces.NewHandler,
	),
	fx.Invoke(func(router *chi.Mux, handler *interfaces.PlayerGameStatisticHandler) {
		router.Post("/log_game_player_statistic/v1/import/{ID}", handler.Handler)

	}),
)
