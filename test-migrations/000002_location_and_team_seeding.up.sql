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
    location_id,
    player_name_encrypted,
    player_name_hash
) 
VALUES
(
    (SELECT id FROM locations WHERE location_name='Pelham Bay'),
    decode('Ex5oYOXNL5Y8m0f7YIqRUA+ZWgjDRZhPV8tOPCGXFtPgE1u75Thp9i0=', 'base64'),-- goofer boofer
    decode('xbdAEiZPI8Bk4uKdBwqswauh25BBd3F5RXJ0AGMivqE=', 'base64')
),
(
    (SELECT id FROM locations WHERE location_name='Elmwood'),
    decode('VrjKjrSWIAAnsA0k1m/e9ItWb8Txm8VtX9MeEULtAAemufa62t8=', 'base64'),-- new player
    decode('UAA5Do+JyUCUuFexLz2aQgiJPcDudRCURW+U0k9S9nE=', 'base64')
),
(
    (SELECT id FROM locations WHERE location_name='Lawrence'),
    decode('Tl5Y+RZYaH16CePouN3Y8sy9s/+QBCvr4x2MJGtIvBk8TtEauLc=', 'base64'),--player one
    decode('s9XALDA7x7nATlrx/hMCYW2gjfoBdFXeitf+EOgULfI=', 'base64')
),
(
    (SELECT id FROM locations WHERE location_name='Lawrence'),
    decode('o/y9bLs16Ta+s/JE2/6f3M1+cuCXFfgoSrEBF8wvxR0lcaFQ1xA=', 'base64'),--player mah
    decode('s/+tK2VVjdvan2ryG33pA8cEpYwrscbJMdF1SY8dWXI=', 'base64')
);

-- ============================
-- Create one Team per Location
-- ============================
INSERT INTO teams (location_id)
SELECT id FROM locations
ON CONFLICT (location_id) DO NOTHING;

-- =============================
-- Seed player saved games
-- =============================
INSERT INTO savedgames (
    location_id,
    winning_player_id,
    winning_player_name_encrypted,
    winning_player_name_hash,
    total_score,
    average_score
)
VALUES
(
    (SELECT id FROM locations WHERE location_name = 'Elmwood'),
    (SELECT id FROM players WHERE player_name_hash = decode('xbdAEiZPI8Bk4uKdBwqswauh25BBd3F5RXJ0AGMivqE=', 'base64')),
    decode('Ex5oYOXNL5Y8m0f7YIqRUA+ZWgjDRZhPV8tOPCGXFtPgE1u75Thp9i0=', 'base64'),-- goofer boofer
    decode('xbdAEiZPI8Bk4uKdBwqswauh25BBd3F5RXJ0AGMivqE=', 'base64'),
    1200,
    400.0
),
(
    (SELECT id FROM locations WHERE location_name = 'Lawrence'),
    (SELECT id FROM players WHERE player_name_hash = decode('UAA5Do+JyUCUuFexLz2aQgiJPcDudRCURW+U0k9S9nE=', 'base64')),
    decode('VrjKjrSWIAAnsA0k1m/e9ItWb8Txm8VtX9MeEULtAAemufa62t8=', 'base64'), -- new player
    decode('UAA5Do+JyUCUuFexLz2aQgiJPcDudRCURW+U0k9S9nE=', 'base64'),
    2200,
    600.25
),
(
    (SELECT id FROM locations WHERE location_name = 'Flushing'),
    (SELECT id FROM players WHERE player_name_hash = decode('s9XALDA7x7nATlrx/hMCYW2gjfoBdFXeitf+EOgULfI=', 'base64')),
    decode('Tl5Y+RZYaH16CePouN3Y8sy9s/+QBCvr4x2MJGtIvBk8TtEauLc=', 'base64'), --player one
    decode('s9XALDA7x7nATlrx/hMCYW2gjfoBdFXeitf+EOgULfI=', 'base64'),
    2200,
    600.25
),
(
    (SELECT id FROM locations WHERE location_name = 'Lawrence'),
    (SELECT id FROM players WHERE player_name_hash = decode('s/+tK2VVjdvan2ryG33pA8cEpYwrscbJMdF1SY8dWXI=', 'base64')),
    decode('o/y9bLs16Ta+s/JE2/6f3M1+cuCXFfgoSrEBF8wvxR0lcaFQ1xA=', 'base64'),--player mah  
    decode('s/+tK2VVjdvan2ryG33pA8cEpYwrscbJMdF1SY8dWXI=', 'base64'),
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