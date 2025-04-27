package player_logs

import (
	"go.uber.org/fx"
	"net/http"
	"skyhawk/internal/player_logs/avg_player_season_statistics"
	"skyhawk/internal/player_logs/avg_team_season_statistics"
	"skyhawk/internal/player_logs/config"
	"skyhawk/internal/player_logs/db"
	"skyhawk/internal/player_logs/player_game_statistic"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			config.LoadConfig,
			db.NewDB,
			db.NewCacheDB,
			NewRouter,
			NewHTTPServer,
		),
		fx.Options(
			player_game_statistic.Module,
			avg_player_season_statistics.Module,
			avg_team_season_statistics.Module,
		),
		fx.Invoke(
			func(*http.Server) {},
		),
	)
}
