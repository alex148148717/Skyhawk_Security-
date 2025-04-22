-- +goose Up

-- Insert team
INSERT INTO teams (id, nba_team_id, name)
VALUES (2, 1610612744, 'Golden State Warriors');

-- Insert players
INSERT INTO players (id, nba_player_id, name)
VALUES (17, 201939, 'Stephen Curry'),
       (18, 203110, 'Draymond Green'),
       (19, 202710, 'Jimmy Butler'),
       (20, 1627741, 'Buddy Hield'),
       (21, 1629645, 'Brandin Podziemski'),
       (22, 1629646, 'Trayce Jackson-Davis'),
       (23, 1629647, 'Moses Moody'),
       (24, 1629648, 'Jonathan Kuminga'),
       (25, 1629649, 'Kevon Looney'),
       (26, 1629650, 'Gary Payton II'),
       (27, 1629651, 'Quinten Post'),
       (28, 1629652, 'Taran Armstrong'),
       (29, 1629653, 'Pat Spencer'),
       (30, 1629654, 'Braxton Key'),
       (31, 1629655, 'Kevin Knox II'),
       (32, 1629656, 'Gui Santos'),
       (33, 4066889, 'Jackson Rowe'),
       (34, 1629658, 'Yuri Collins'),
       (35, 1629659, 'Reece Beekman'),
       (36, 1629660, 'Anthony Melton'),
       (37, 1629661, 'Kyle Anderson'),
       (38, 1629662, 'Lindy Waters III')
    ON CONFLICT (nba_player_id) DO NOTHING;

-- Assign jersey numbers for season 2024
INSERT INTO player_team_numbers (player_id, team_id, season_year, jersey_number)
VALUES (17, 2, 2024, 30),
       (18, 2, 2024, 23),
       (19, 2, 2024, 22),
       (20, 2, 2024, 7),
       (21, 2, 2024, 2),
       (22, 2, 2024, 32),
       (23, 2, 2024, 4),
       (24, 2, 2024, 00),
       (25, 2, 2024, 5),
       (26, 2, 2024, 0),
       (27, 2, 2024, 21),
       (28, 2, 2024, 1),
       (29, 2, 2024, 3),
       (30, 2, 2024, 12),
       (31, 2, 2024, 31),
       (32, 2, 2024, 15),
       (33, 2, 2024, 14),
       (34, 2, 2024, 11),
       (35, 2, 2024, 6),
       (36, 2, 2024, 8),
       (37, 2, 2024, 5),
       (38, 2, 2024, 13);

-- Assign jersey numbers for season 2025
INSERT INTO player_team_numbers (player_id, team_id, season_year, jersey_number)
VALUES (17, 2, 2025, 30),
       (18, 2, 2025, 23),
       (19, 2, 2025, 22),
       (20, 2, 2025, 7),
       (21, 2, 2025, 2),
       (22, 2, 2025, 32),
       (23, 2, 2025, 4),
       (24, 2, 2025, 00),
       (25, 2, 2025, 5),
       (26, 2, 2025, 0),
       (27, 2, 2025, 21),
       (28, 2, 2025, 1),
       (29, 2, 2025, 3),
       (30, 2, 2025, 12),
       (31, 2, 2025, 31),
       (32, 2, 2025, 15),
       (33, 2, 2025, 14),
       (34, 2, 2025, 11),
       (35, 2, 2025, 6),
       (36, 2, 2025, 8),
       (37, 2, 2025, 5),
       (38, 2, 2025, 13);

-- +goose Down

-- Remove player-team numbers
DELETE
FROM player_team_numbers
WHERE player_id BETWEEN 17 AND 38;

-- Remove players
DELETE
FROM players
WHERE id BETWEEN 17 AND 38;

-- Remove team
DELETE
FROM teams
WHERE id = 2;
