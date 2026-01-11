-- ============================
-- Locations
-- ============================
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    location_name VARCHAR(30) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================
-- Teams
-- ============================
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    location_id INT NOT NULL REFERENCES locations(id) ON DELETE CASCADE,

    UNIQUE (location_id, id)
);

-- ============================
-- Players
-- ============================
CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    player_name VARCHAR(60) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    location_id INT NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    team_id INT REFERENCES teams(id) ON DELETE SET NULL
);

-- ============================
-- Saved Games
-- ============================
CREATE TABLE saved_games (
    id INT PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    winning_player_name VARCHAR(60),
    total_score INT NOT NULL,
    average_score DOUBLE PRECISION NOT NULL

    location_id INT NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    winning_team_id INT REFERENCES teams(id) ON DELETE SET NULL,
    winning_player_id INT REFERENCES players(id) ON DELETE SET NULL,
);

-- ============================
-- Saved Game Teams
-- ============================
CREATE TABLE saved_game_teams (
    id INT PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    team_score INT NOT NULL,

    team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    saved_game_id INT NOT NULL REFERENCES saved_games(id) ON DELETE CASCADE,

    UNIQUE (saved_game_id, team_id)
);

-- ============================
-- Saved Game Players
-- ============================
CREATE TABLE saved_game_players (
    id INT PRIMARY KEY,
    player_name VARCHAR(60) NOT NULL,
    player_score INT NOT NULL,

    saved_game_id INT NOT NULL REFERENCES saved_games(id) ON DELETE CASCADE,
    player_id INT REFERENCES players(id) ON DELETE SET NULL,

    UNIQUE (saved_game_id, player_id)
);