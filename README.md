# Player Stats Service

This is a Go-based backend service for managing and retrieving NBA player statistics.  
It uses PostgreSQL for persistent storage and optionally supports AWS DynamoDB (with optional DAX support) for caching.

---

## Technologies Used

- Go (Golang)
- PostgreSQL
- DynamoDB (with optional DAX support)

## Why PostgreSQL?

- Managed relational database available in all major cloud platforms
- Suitable for structured and aggregated data
- Scales well for expected data volume
- Simplifies reporting and queries using SQL

## Why DynamoDB?

- Managed NoSQL service, integrated into AWS
- Effective for caching player pages or frequently accessed data
- Requires no infrastructure management
- Optional DAX integration for higher performance under heavy load
- Fixed-cost pricing option compared to compute-based alternatives

## Runtime Configuration

All configuration options are passed via command-line flags:

```bash
./app 
```

## Example URLs

To query aggregated stats:

- **Team stats** (NBA team ID: `1610612747`, season: `2024`):  
  `http://localhost:8081/log_game_player_statistic/v1/season/2024/team/1610612747`

- **Player stats** (NBA player ID: `2544`, season: `2024`):  
  `http://localhost:8081/log_game_player_statistic/v1/season/2024/player/2544`

> Note: These endpoints will return empty until data is ingested.

## Importing Events for Game 1

To import player statistics for a specific game (e.g., Lakers vs Warriors, season 2024), use the following command:

```bash
UUID=$(uuidgen)
echo "Using UUID: $UUID"

curl -X POST "http://localhost:8081/log_game_player_statistic/v1/import/$UUID" \
  -H "Content-Type: application/json" \
  --data-binary @examples/game_2024_lakers_vs_warriors_stats.json
```
```bash
UUID=$(uuidgen)
echo "Using UUID: $UUID"

curl -X POST "http://localhost:8081/log_game_player_statistic/v1/import/$UUID" \
  -H "Content-Type: application/json" \
  --data-binary @examples/game_2024_lakers_vs_warriors_stats_game_2.json
```
```bash
UUID=$(uuidgen)
echo "Using UUID: $UUID"

curl -X POST "http://localhost:8081/log_game_player_statistic/v1/import/$UUID" \
  -H "Content-Type: application/json" \
  --data-binary @examples/game_2024_lakers_vs_warriors_stats_game_3.json
```

## Example URLs

The following endpoints return **aggregated statistics** based on imported game data:

- **Team stats** (NBA team ID: `1610612747`, season: `2024`):  
  `http://localhost:8081/log_game_player_statistic/v1/season/2024/team/1610612747`

- **Player stats** (NBA player ID: `2544`, season: `2024`):  
  `http://localhost:8081/log_game_player_statistic/v1/season/2024/player/2544`

> These endpoints now return **aggregated data** after importing game events.


## Import Format

Each imported payload must be a JSON array of player statistics objects.  
Example:

```json
[
  {
    "id": 38,
    "season_year": 2024,
    "game_id": 1,
    "team_id": 2,
    "player_id": 38,
    "points": 4,
    "rebounds": 8,
    "assists": 0,
    "steals": 2,
    "blocks": 1,
    "fouls": 3,
    "turnovers": 2,
    "minutes_played": 11.27
  }
]
```

## Database Schema

This service uses a PostgreSQL database to store player and team information, jersey numbers, and game statistics.

### Tables and Relationships

#### `players`

Stores player records.

- `id` (primary key)
- `nba_player_id` – unique identifier from the NBA
- `name` – full name of the player

#### `teams`

Stores team records.

- `id` (primary key)
- `nba_team_id` – unique identifier from the NBA
- `name` – team name

#### `player_team_numbers`

Links players to teams with jersey numbers for a given season.

- `id` (primary key)
- `player_id` (foreign key → `players.id`)
- `team_id` (foreign key → `teams.id`)
- `season_year` – e.g., 2024
- `jersey_number` – number worn that season

#### `player_stats_raw`

Stores per-game player statistics.

- `id` (primary key)
- `job_id` – import/job UUID
- `season_id` – year of the season
- `game_id` – local game identifier
- `player_id` (foreign key → `players.id`)
- `team_id` (foreign key → `teams.id`)
- `points`, `rebounds`, `assists`, `steals`, `blocks`, `fouls`, `turnovers`, `minutes_played`

### Relationships

- A player can play for multiple teams in different seasons (`players` ⇄ `teams` via `player_team_numbers`)
- A player can have multiple stat entries across games (`player_stats_raw`)
- Each stat record links to both the player and the team involved in the game

## Data Flow

1. **Receive Update/Sync Request**
   The service receives a request (e.g., HTTP POST) to update or synchronize player game data into the database (PostgreSQL and/or DynamoDB).
2. **Sync Data to Database**
   The data is inserted or updated in the PostgreSQL `player_stats_raw` table and optionally cached in DynamoDB (or DAX if enabled).
3. **Generate and Cache HTML Page**
   After the database update, the service regenerates the relevant HTML page based on the requested keys (e.g., team or player for a specific season).
4. **Store in Cache**
   The generated HTML page is stored in the cache (DynamoDB) using a defined key pattern for fast retrieval.
5. **Serve Page from Cache**
   When a user requests a team or player page, the service fetches the pre-rendered HTML content from the cache and returns it as the response.
6. 