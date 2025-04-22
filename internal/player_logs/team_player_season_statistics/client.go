package team_player_season_statistics

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"

	"fmt"
	"html/template"
)

//go:embed template.html
var pageTemplateByte string
var pageTpl *template.Template

func init() {
	var err error
	pageTpl, err = template.New("teamPlayerPage").Parse(pageTemplateByte)
	if err != nil {
		panic(fmt.Errorf("failed to parse template: %w", err))
	}
}

type AveragePlayerSeasonClient interface {
	AveragePlayerSeason(ctx context.Context, ids []int32) ([]AveragePlayerSeason, error)
}
type CacheClient interface {
	PutItem(ctx context.Context, key string, value []byte) error
	GetItem(ctx context.Context, key string) ([]byte, error)
}

type Client struct {
	averagePlayerSeasonClient AveragePlayerSeasonClient
	cacheClient               CacheClient
}

func New(r chi.Router, averagePlayerSeasonClient AveragePlayerSeasonClient, cacheClient CacheClient) *Client {

	c := &Client{
		averagePlayerSeasonClient: averagePlayerSeasonClient,
		cacheClient:               cacheClient,
	}
	r.Get("/season/{SeasonID}/player/{NbaPlayerID}", c.SeasonTeamPlayerStatisticsHandler)
	return c

}
func keyGenerate(seasonID, playerID int) string {
	return fmt.Sprintf("AveragePlayerSeasonClient_%d_%d", seasonID, playerID)

}

func (c *Client) Sync(ctx context.Context, ids []int32) error {
	averagePlayerSeason, err := c.averagePlayerSeasonClient.AveragePlayerSeason(ctx, ids)
	if err != nil {
		return err
	}
	for _, a := range averagePlayerSeason {
		key := keyGenerate(a.SeasonID, a.NbaPlayerID)
		var buf bytes.Buffer
		if err := pageTpl.Execute(&buf, a); err != nil {
			return err
		}
		data := buf.Bytes()

		c.cacheClient.PutItem(ctx, key, data)
	}
	_ = averagePlayerSeason
	return nil
}

func (c *Client) SeasonTeamPlayerStatisticsHandler(w http.ResponseWriter, r *http.Request) {
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
	key := keyGenerate(seasonID, nbaPlayerID)
	fmt.Printf("key %s\n", key)
	html, err := c.cacheClient.GetItem(ctx, key)
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
	w.Write(html)
	return

}
