-- Delete teams tied to seeded locations
DELETE FROM teams
WHERE location_id IN (
    SELECT id FROM locations
    WHERE location_name IN (
        'Elmwood',
        'Pelham Bay',
        'Grand Concourse',
        'W 154th St',
        '5030 Broadway',
        'Flushing',
        'Lawrence',
        'Port Richmond'
    )
);

-- Delete the seeded locations
DELETE FROM locations
WHERE location_name IN (
    'Elmwood',
    'Pelham Bay',
    'Grand Concourse',
    'W 154th St',
    '5030 Broadway',
    'Flushing',
    'Lawrence',
    'Port Richmond'
);