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

INSERT INTO players (player_name, location_id) VALUES('playerone', (SELECT id FROM locations WHERE location_name='Pelham Bay'));
INSERT INTO players (player_name, location_id) VALUES('playertwo', (SELECT id FROM locations WHERE location_name='Pelham Bay'));
INSERT INTO players (player_name, location_id) VALUES('playerthree', (SELECT id FROM locations WHERE location_name='Elmwood'));
INSERT INTO players (player_name, location_id) VALUES('playerfour', (SELECT id FROM locations WHERE location_name='Elmwood'));

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
)
ON CONFLICT (savedgameid) DO NOTHING;
