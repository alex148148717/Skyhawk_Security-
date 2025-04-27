package player_logs

import (
	"go.uber.org/fx"
	"net/http"
	"skyhawk/internal/player_logs/avg_player_season_statistics"
	"skyhawk/internal/player_logs/avg_team_season_statistics"
	"skyhawk/internal/player_logs/config"
	"skyhawk/internal/player_logs/db"
	"skyhawk/internal/player_logs/player_game_statistic"
	"testing"
)

func LoadConfig() *config.Config {
	return &config.Config{
		AppEnv:         "development",
		DatabaseURL:    "host=localhost user=user password=pass dbname=player_stats port=5444 sslmode=disable TimeZone=UTC",
		ServerPort:     "8081",
		CacheTableName: "cache",
		Dynamodb: config.Dynamodb{
			Region:       "us-west-2",
			Endpoint:     "http://localhost:8000",
			DaxHostPorts: []string{},
			UseDax:       false,
		},
	}
}

func TestApp(t *testing.T) {

	fx.New(
		fx.Provide(
			LoadConfig,
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
	).Run()

}
