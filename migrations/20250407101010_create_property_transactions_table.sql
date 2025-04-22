-- +goose Up
-- +goose StatementBegin
CREATE TABLE player_stats_raw (
    id              INT PRIMARY KEY,
    job_id          UUID,
    season_id       INT,
    game_id         INT,
    player_id INT NOT NULL REFERENCES players(id),
    team_id INT NOT NULL REFERENCES teams(id),
    points          SMALLINT,
    rebounds        SMALLINT,
    assists         SMALLINT,
    steals          SMALLINT,
    blocks          SMALLINT,
    fouls           SMALLINT,
    turnovers       SMALLINT,
    minutes_played  REAL
);
CREATE INDEX idx_job_season_team ON player_stats_raw (job_id, season_id, team_id);
CREATE INDEX idx_season_team ON player_stats_raw (season_id, team_id);
CREATE INDEX idx_season_team_player ON player_stats_raw (season_id, team_id, player_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_season_team_player;
DROP INDEX IF EXISTS idx_season_team;
DROP INDEX IF EXISTS idx_job_season_team;
DROP TABLE IF EXISTS player_stats_raw;
-- +goose StatementEnd
