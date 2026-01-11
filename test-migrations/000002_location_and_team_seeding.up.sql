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
INSERT INTO teams (location_id)
SELECT l.id
FROM locations l
WHERE l.location_name IN (
    'Elmwood',
    'Pelham Bay',
    'Grand Concourse',
    'W 154th St',
    '5030 Broadway',
    'Flushing',
    'Lawrence',
    'Port Richmond'
)
AND NOT EXISTS (
    SELECT 1
    FROM teams t
    WHERE t.location_id = l.id
);