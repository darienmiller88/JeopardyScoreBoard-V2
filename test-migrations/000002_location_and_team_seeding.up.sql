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