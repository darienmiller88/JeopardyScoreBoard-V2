CREATE TABLE IF NOT EXISTS SavedGamesTeams (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- table specific fields
    team_id INTEGER NOT NULL,
    saved_game_id INTEGER NOT NULL,
    team_score INTEGER,
    
    -- fk constraints
    FOREIGN KEY (team_id) REFERENCES Teams(id),
    FOREIGN KEY (saved_game_id) REFERENCES SavedGames(id) ON DELETE CASCADE,
    UNIQUE(team_id, saved_game_id)
);

CREATE TABLE IF NOT EXISTS SavedGamesPlayers (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- table specific fields
    player_id INTEGER NOT NULL,
    saved_game_id INTEGER NOT NULL,
    player_score INTEGER,
    
    -- fk constraints
    FOREIGN KEY (player_id) REFERENCES Players(id),
    FOREIGN KEY (saved_game_id) REFERENCES SavedGames(id) ON DELETE CASCADE,
    UNIQUE(player_id, saved_game_id)
);

DROP TABLE IF EXISTS GameParticipation;