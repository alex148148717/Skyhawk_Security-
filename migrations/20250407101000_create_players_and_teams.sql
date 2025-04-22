-- +goose Up
-- +goose StatementBegin
CREATE TABLE players (
    id INT PRIMARY KEY,
    nba_player_id INTEGER NOT NULL UNIQUE,
    name TEXT NOT NULL
);

CREATE TABLE teams (
    id INT PRIMARY KEY,
    nba_team_id INT NOT NULL UNIQUE,
    name TEXT NOT NULL
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
