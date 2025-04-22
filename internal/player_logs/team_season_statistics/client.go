package team_season_statistics

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
)

//go:embed template.html
var pageTemplateByte string
var pageTpl *template.Template

func init() {
	var err error
	pageTpl, err = template.New("teamPage").Parse(pageTemplateByte)
	if err != nil {
		panic(fmt.Errorf("failed to parse template: %w", err))
	}
}

type AverageTeamSeasonClient interface {
	AverageTeamSeason(ctx context.Context, ids []int32) ([]AverageTeamSeason, error)
}
type CacheClient interface {
	PutItem(ctx context.Context, key string, value []byte) error
	GetItem(ctx context.Context, key string) ([]byte, error)
}

type Client struct {
	averageTeamSeasonClient AverageTeamSeasonClient
	cacheClient             CacheClient
}

func New(r chi.Router, averageTeamSeasonClient AverageTeamSeasonClient, cacheClient CacheClient) *Client {
	c := &Client{
		averageTeamSeasonClient: averageTeamSeasonClient,
		cacheClient:             cacheClient,
	}
	r.Get("/season/{SeasonID}/team/{NbaTeamID}", c.SeasonTeamStatisticsHandler)

	return c
}

func (c *Client) Sync(ctx context.Context, ids []int32) error {
	averageTeamSeason, err := c.averageTeamSeasonClient.AverageTeamSeason(ctx, ids)
	if err != nil {
		return err
	}

	for _, a := range averageTeamSeason {
		key := keyGenerate(a.SeasonID, a.NbaTeamID)
		var buf bytes.Buffer
		if err := pageTpl.Execute(&buf, a); err != nil {
			return err
		}
		data := buf.Bytes()
		c.cacheClient.PutItem(ctx, key, data)
	}
	_ = averageTeamSeason
	return nil
}
func keyGenerate(seasonID, teamID int) string {
	return fmt.Sprintf("AverageTeamSeasonClient_%d_%d", seasonID, teamID)

}

func (c *Client) SeasonTeamStatisticsHandler(w http.ResponseWriter, r *http.Request) {
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
	key := keyGenerate(seasonID, nbaTeamID)
	html, err := c.cacheClient.GetItem(ctx, key)
	if err != nil {
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
