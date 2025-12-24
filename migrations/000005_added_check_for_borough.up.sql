ALTER TABLE locations
ADD CONSTRAINT borough_check CHECK (borough IN ('Brooklyn', 'Bronx', 'Manhattan', 'Staten Island', 'Queens'));