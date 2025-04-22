-- +goose Up

-- Insert team (Lakers)
INSERT INTO teams (id, nba_team_id, name)
VALUES (1, 1610612747, 'Los Angeles Lakers');

-- Insert all Lakers players for 2024 season
INSERT INTO players (id, nba_player_id, name)
VALUES (1, 2544, 'LeBron James'),
       (2, 203076, 'Anthony Davis'),
       (3, 1627746, 'Angelo Russell'),
(4, 1628398, 'Austin Reaves'),
(5, 1630553, 'Rui Hachimura'),
(6, 1626179, 'Jarred Vanderbilt'),
(7, 1628366, 'Gabe Vincent'),
(8, 1627826, 'Taurean Prince'),
(9, 1631262, 'Max Christie'),
(10, 1630581, 'Jaxson Hayes'),
(11, 1629013, 'Cam Reddish'),
(12, 1631121, 'Colin Castleton'),
(13, 1631246, 'Alex Fudge'),
(14, 1631293, 'Moi Hodge'),
       (15, 1631122, 'Maxwell Lewis'),
       (16, 1626174, 'Christian Wood');

-- Assign jersey numbers for season 2024
INSERT INTO player_team_numbers (player_id, team_id, season_year, jersey_number)
VALUES (1, 1, 2024, 6),
       (2, 1, 2024, 3),
       (3, 1, 2024, 1),
       (4, 1, 2024, 15),
       (5, 1, 2024, 28),
       (6, 1, 2024, 2),
       (7, 1, 2024, 7),
       (8, 1, 2024, 12),
       (9, 1, 2024, 10),
       (10, 1, 2024, 11),
       (11, 1, 2024, 5),
       (12, 1, 2024, 14),
       (13, 1, 2024, 17),
       (14, 1, 2024, 55),
       (15, 1, 2024, 21),
       (16, 1, 2024, 35);

-- +goose Down

DELETE
FROM player_team_numbers
WHERE player_id BETWEEN 1 AND 16;
DELETE
FROM players
WHERE id BETWEEN 1 AND 16;
DELETE
FROM teams
WHERE id = 1;
