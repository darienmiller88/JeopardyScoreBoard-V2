ALTER TABLE locations ADD COLUMN borough VARCHAR(30) NOT NULL 
CONSTRAINT borough_check CHECK (borough IN ('Brooklyn', 'Bronx', 'Manhattan', 'Staten Island', 'Queens'));