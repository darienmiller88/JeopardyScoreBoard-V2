-- ============================
-- Locations
-- ============================
CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    location_name VARCHAR(30) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================
-- Teams
-- ============================
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    location_id INT NOT NULL REFERENCES locations(id) ON DELETE CASCADE,

    UNIQUE (location_id)
);

-- ============================
-- Players
-- ============================
CREATE TABLE IF NOT EXISTS players (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    player_name_encrypted BYTEA NOT NULL,
    player_name_hash BYTEA NOT NULL,

    location_id INT NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    team_id INT REFERENCES teams(id) ON DELETE SET NULL,

    UNIQUE(player_name_hash, location_id)
);

-- ============================
-- Saved Games
-- ============================
CREATE TABLE IF NOT EXISTS savedgames (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    total_score INT NOT NULL,
    average_score DOUBLE PRECISION NOT NULL,
    
    winning_player_name_encrypted BYTEA,
    winning_player_name_hash BYTEA,

    location_id INT NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    winning_team_id INT REFERENCES teams(id) ON DELETE SET NULL,
    winning_player_id INT REFERENCES players(id) ON DELETE SET NULL
);


-- ============================
-- Saved Game Teams
-- ============================
CREATE TABLE IF NOT EXISTS savedgamesteams (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    team_score INT NOT NULL,

    team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    saved_game_id INT NOT NULL REFERENCES savedgames(id) ON DELETE CASCADE,

    UNIQUE (saved_game_id, team_id)
);

-- ============================
-- Saved Game Players
-- ============================
CREATE TABLE IF NOT EXISTS savedgamesplayers (
    id SERIAL PRIMARY KEY,
    player_score INT NOT NULL,

    player_name_encrypted BYTEA NOT NULL,
    player_name_hash BYTEA NOT NULL,

    saved_game_id INT NOT NULL REFERENCES savedgames(id) ON DELETE CASCADE,
    player_id INT REFERENCES players(id) ON DELETE SET NULL,

    UNIQUE (saved_game_id, player_id)
);

