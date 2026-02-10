-- ============================
-- Seed Locations
-- ============================
INSERT INTO locations (location_name)
VALUES
    ('Elmwood'),
    ('Pelham Bay'),
    ('Grand Concourse'),
    ('W 154th St'),
    ('5030 Broadway'),
    ('Flushing'),
    ('Lawrence'),
    ('Port Richmond')
ON CONFLICT (location_name) DO NOTHING;


-- ============================
-- Seed Players
-- ============================
INSERT INTO players (
    location_id
    player_name_encrypted,
    player_name_hash
) 
VALUES
(
    (SELECT id FROM locations WHERE location_name='Pelham Bay')
    decode('Ex5oYOXNL5Y8m0f7YIqRUA+ZWgjDRZhPV8tOPCGXFtPgE1u75Thp9i0=', 'base64'),-- goofer boofer
    decode('xbdAEiZPI8Bk4uKdBwqswauh25BBd3F5RXJ0AGMivqE=', 'base64'),
),
(
    (SELECT id FROM locations WHERE location_name='Elmwood')
    decode('VrjKjrSWIAAnsA0k1m/e9ItWb8Txm8VtX9MeEULtAAemufa62t8=', 'base64'),-- new player
    decode('UAA5Do+JyUCUuFexLz2aQgiJPcDudRCURW+U0k9S9nE=', 'base64'),
),
(
    (SELECT id FROM locations WHERE location_name='Lawrence')
    decode('Tl5Y+RZYaH16CePouN3Y8sy9s/+QBCvr4x2MJGtIvBk8TtEauLc=', 'base64'),--player one
    decode('wVwk5c+ad0t3B5YbaKY3FsrVpO/PWm7vZFDBlLW/j50=', 'base64'),
),
(
    (SELECT id FROM locations WHERE location_name='Lawrence')
    decode('o/y9bLs16Ta+s/JE2/6f3M1+cuCXFfgoSrEBF8wvxR0lcaFQ1xA=', 'base64'),--player mah
    decode('wVwk5c+ad0t3B5YbaKY3FsrVpO/PWm7vZFDBlLW/j50=', 'base64'),  
);

-- ============================
-- Create one Team per Location
-- ============================
INSERT INTO teams (location_id) VALUES((SELECT id FROM locations WHERE location_name='Elmwood'));
INSERT INTO teams (location_id) VALUES((SELECT id FROM locations WHERE location_name='Lawrence'));
INSERT INTO teams (location_id) VALUES((SELECT id FROM locations WHERE location_name='Flushing'));
INSERT INTO teams (location_id) VALUES((SELECT id FROM locations WHERE location_name='Grand Concourse'));
INSERT INTO teams (location_id) VALUES((SELECT id FROM locations WHERE location_name='Port Richmond'));
INSERT INTO teams (location_id) VALUES((SELECT id FROM locations WHERE location_name='5030 Broadway'));
INSERT INTO teams (location_id) VALUES((SELECT id FROM locations WHERE location_name='W 154th St'));
INSERT INTO teams (location_id) VALUES((SELECT id FROM locations WHERE location_name='Pelham Bay'));


-- =============================
-- Seed player saved games
-- =============================
INSERT INTO savedgames (
    location_id,
    winning_player_id,
    winning_player_name,
    total_score,
    average_score
)
VALUES
(
    (SELECT id FROM locations WHERE location_name = 'Elmwood'),
    (SELECT id FROM players WHERE player_name = 'playerone'),
    'playerone',
    1200,
    400.0
),
(
    (SELECT id FROM locations WHERE location_name = 'Lawrence'),
    (SELECT id FROM players WHERE player_name = 'playerone'),
    'playerone',
    2200,
    600.25
),
(
    (SELECT id FROM locations WHERE location_name = 'Flushing'),
    (SELECT id FROM players WHERE player_name = 'playertwo'),
    'playertwo',
    2200,
    600.25
),
(
    (SELECT id FROM locations WHERE location_name = 'Lawrence'),
    (SELECT id FROM players WHERE player_name = 'playerthree'),
    'playerthree',
    2200,
    600.25
);
-- Removed: ON CONFLICT (savedgameid) DO NOTHING;
-- This would only work if savedgameid is defined as a unique constraint



-- =============================
-- Seed team saved games
-- =============================
INSERT INTO savedgames (
    location_id,
    winning_team_id, 
    total_score, 
    average_score
)
VALUES
(
    (SELECT id FROM locations WHERE location_name = 'Elmwood'),
    1,
    1200,
    400.0
),
(
    (SELECT id FROM locations WHERE location_name = 'Lawrence'),
    2,
    2200,
    600.25
),
(
    (SELECT id FROM locations WHERE location_name = 'Flushing'),
    3,
    2200,
    600.25
),
(
    (SELECT id FROM locations WHERE location_name = 'Lawrence'),
    4,
    2200,
    600.25
);
-- Removed: ON CONFLICT (savedgameid) DO NOTHING;