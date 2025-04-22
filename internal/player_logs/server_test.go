package player_logs

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"skyhawk/internal/player_logs/player_game_statistic"
	"skyhawk/internal/player_logs/player_log_cache"
	"skyhawk/internal/player_logs/team_player_season_statistics"
	"skyhawk/internal/player_logs/team_season_statistics"
	"testing"
)

//go:embed mockdata.json
var sampleJSON []byte

func TestServer(t *testing.T) {
	ctx := context.Background()
	_ = ctx
	dsn := "host=localhost user=user password=pass dbname=player_stats port=5444 sslmode=disable TimeZone=UTC"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	playerGameStatisticDbClient := player_game_statistic.NewDBClient(db)
	teamPlayerSeasonStatisticsDbClient := team_player_season_statistics.NewDBClient(db)
	teamSeasonStatisticsDbClient := team_season_statistics.NewDBClient(db)

	cfg := &aws.Config{
		Region: aws.String("us-west-2"),
	}

	cfg.Endpoint = aws.String("http://localhost:8000")
	cfg.Credentials = credentials.NewStaticCredentials("fake", "fake", "")

	sess := session.Must(session.NewSession(cfg))

	svc := dynamodb.New(sess)

	playerLogCacheClient := player_log_cache.New("cache", svc)

	r := chi.NewRouter()
	r.Route("/log_game_player_statistic/v1/", func(r chi.Router) {

		teamSeasonStatisticsClient := team_season_statistics.New(r, teamSeasonStatisticsDbClient, playerLogCacheClient)
		teamPlayerSeasonStatistics := team_player_season_statistics.New(r, teamPlayerSeasonStatisticsDbClient, playerLogCacheClient)
		syncClients := []player_game_statistic.SyncClient{
			teamSeasonStatisticsClient,
			teamPlayerSeasonStatistics,
		}

		player_game_statistic.New(r, playerGameStatisticDbClient, syncClients)

	})

	httpServer := httptest.NewServer(r)
	defer httpServer.Close()
	//import data
	id := uuid.NewString()
	fmt.Printf("job id %s\n", id)
	importUrl := fmt.Sprintf("%s/log_game_player_statistic/v1/import/%s", httpServer.URL, id)
	resp, err := http.Post(importUrl, "application/json", bytes.NewBuffer(sampleJSON))
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}
	defer resp.Body.Close()
	//get data
	seasonID := 2022
	teamID := 1
	seasonTeamUrl := fmt.Sprintf("%s/log_game_player_statistic/v1/season/%d/player/%d", httpServer.URL, seasonID, teamID)

	resp2, err := http.Get(seasonTeamUrl)
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}
	defer resp2.Body.Close()
	bodyBytes, err := io.ReadAll(resp2.Body)
	if err != nil {
		log.Fatal("Failed to read body:", err)
	}

	fmt.Println("📄 Response body:")
	fmt.Println(string(bodyBytes))
}

type SeasonAverage struct {
	JobID         string
	SeasonID      uint64
	TeamID        uint64
	AvgPoints     float32
	Rebounds      float32
	Assists       float32
	Steals        float32
	Blocks        float32
	Fouls         float32
	Turnovers     float32
	MinutesPlayed float32
}
