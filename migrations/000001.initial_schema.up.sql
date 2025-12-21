CREATE TABLE IF NOT EXISTS Locations(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- field name
    location_name VARCHAR(60)
);

CREATE TABLE IF NOT EXISTS Teams(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign key
    location_id UUID NOT NULL UNIQUE

    -- constraint
    FOREIGN KEY (location_id) REFERENCES Locations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS SavedGames(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- table specific fields
    winner VARCHAR(60) NOT NULL,
    total_score INT NOT NULL,
    average_score DOUBLE PRECISION NOT NULL,
    is_team_game BOOLEAN NOT NULL,
 
    -- foreign key 
    location_id UUID NOT NULL,

    --Constraint
    FOREIGN KEY (location_id) REFERENCES Locations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Players(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Table specific fields
    player_name varchar(60) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    left_at TIMESTAMPTZ,

    -- Foreign keys
    location_id UUID NOT NULL,
    team_id UUID,

    -- Constraints
    FOREIGN KEY (team_id, location_id) REFERENCES Teams(id, location_id)
);

CREATE TABLE IF NOT EXISTS GameParticipation(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Table specific fields
    saved_game_id UUID NOT NULL,
    player_id UUID NOT NULL,
    team_id UUID,

    -- Constraints
    FOREIGN KEY (saved_game_id) REFERENCES SavedGames(id) ON DELETE CASCADE
    FOREIGN KEY (player_id) REFERENCES Players(id)
    FOREIGN KEY (team_id) REFERENCES Teams(id)
);