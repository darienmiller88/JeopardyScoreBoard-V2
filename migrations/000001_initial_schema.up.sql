CREATE TABLE IF NOT EXISTS Locations(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- field name
    location_name VARCHAR(60) 
);

CREATE TABLE IF NOT EXISTS Teams(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign key (UNIQUE ensures one team per location)
    location_id SERIAL NOT NULL UNIQUE,

    -- constraint
    FOREIGN KEY (location_id) REFERENCES Locations(id) ON DELETE CASCADE,
    
    -- Explicit composite unique constraint for the foreign key reference
    UNIQUE (id, location_id)
);

CREATE TABLE IF NOT EXISTS SavedGames(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- table specific fields
    winner VARCHAR(60) NOT NULL,
    total_score INT NOT NULL,
    average_score DECIMAL NOT NULL,
    is_team_game BOOLEAN NOT NULL,
 
    -- foreign key 
    location_id SERIAL NOT NULL,

    -- Constraint
    FOREIGN KEY (location_id) REFERENCES Locations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Players(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Table specific fields
    player_name VARCHAR(60) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    left_at TIMESTAMPTZ,

    -- Foreign keys
    location_id SERIAL NOT NULL,
    team_id SERIAL,

    -- Constraints
    FOREIGN KEY (location_id) REFERENCES Locations(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id, location_id) REFERENCES Teams(id, location_id)
);

CREATE TABLE IF NOT EXISTS GameParticipation(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Table specific fields
    saved_game_id SERIAL NOT NULL,
    player_id SERIAL NOT NULL,
    team_id SERIAL,

    -- Constraints
    FOREIGN KEY (saved_game_id) REFERENCES SavedGames(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES Players(id),
    FOREIGN KEY (team_id) REFERENCES Teams(id),
    
    -- Prevent duplicate participations
    UNIQUE (saved_game_id, player_id)
);