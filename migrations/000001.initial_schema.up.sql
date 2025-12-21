CREATE TABLE IF NOT EXISTS Locations(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    location_name VARCHAR(60)
);

CREATE TABLE IF NOT EXISTS Teams(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    location_id UUID NOT NULL UNIQUE

    FOREIGN KEY (location_id) REFERENCES Locations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS SavedGames(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    location_id UUID NOT NULL,
    winner VARCHAR(60) NOT NULL,
    total_score INT NOT NULL,
    average_score DOUBLE PRECISION NOT NULL,
    is_team_game BOOLEAN NOT NULL,
 
    FOREIGN KEY (location_id) REFERENCES Locations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Players(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    player_name varchar(60) NOT NULL,
    location_id UUID NOT NULL,
    team_id UUID,

    FOREIGN KEY (team_id, location_id) REFERENCES Teams(id, location_id)
);