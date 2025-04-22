-- +goose Up
-- +goose StatementBegin



CREATE TABLE player_team_numbers (
    id SERIAL PRIMARY KEY,
    player_id INT  NOT NULL REFERENCES players(id),
    team_id INT  NOT NULL REFERENCES teams(id),
    season_year INT NOT NULL,
    jersey_number INT NOT NULL
);
CREATE INDEX idx_player_team_numbers ON player_team_numbers (season_year, team_id, player_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS player_team_numbers;
-- +goose StatementEnd
