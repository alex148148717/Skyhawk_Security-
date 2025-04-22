package player_logs

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-dax-go/dax"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"skyhawk/internal/player_logs/player_game_statistic"
	"skyhawk/internal/player_logs/player_log_cache"
	"skyhawk/internal/player_logs/team_player_season_statistics"
	"skyhawk/internal/player_logs/team_season_statistics"
	"syscall"
	"time"

	"database/sql"
)

func New(ctx context.Context, port int, dsn string, region string, endpoint string, tableName string, isDax bool, daxHostPorts []string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	playerGameStatisticDbClient := player_game_statistic.NewDBClient(db)
	teamPlayerSeasonStatisticsDbClient := team_player_season_statistics.NewDBClient(db)
	teamSeasonStatisticsDbClient := team_season_statistics.NewDBClient(db)

	cfg := &aws.Config{
		Region: aws.String(region),
	}
	if endpoint != "" {
		cfg.Endpoint = aws.String(endpoint)
		cfg.Credentials = credentials.NewStaticCredentials("fake", "fake", "")
	}

	sess := session.Must(session.NewSession(cfg))
	var svc dynamodbiface.DynamoDBAPI
	if isDax {
		daxCfg := dax.DefaultConfig()
		daxCfg.HostPorts = daxHostPorts
		daxCfg.Region = *cfg.Region

		svc, err = dax.New(daxCfg)
		if err != nil {
			return err
		}

	} else {
		svc = dynamodb.New(sess)

	}

	playerLogCacheClient := player_log_cache.New(tableName, svc)

	r := chi.NewRouter()
	r.Route("/log_game_player_statistic/v1/", func(r chi.Router) {

		teamSeasonStatisticsClient := team_season_statistics.New(r, teamSeasonStatisticsDbClient, playerLogCacheClient)
		teamPlayerSeasonStatistics := team_player_season_statistics.New(r, teamPlayerSeasonStatisticsDbClient, playerLogCacheClient)
		syncClients := []player_game_statistic.SyncClient{teamSeasonStatisticsClient, teamPlayerSeasonStatistics}

		player_game_statistic.New(r, playerGameStatisticDbClient, syncClients)

	})
	// --- Server ---
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	// Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Shutting down server gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}

	<-idleConnsClosed
	log.Println("Server shutdown complete")
	return nil

}
