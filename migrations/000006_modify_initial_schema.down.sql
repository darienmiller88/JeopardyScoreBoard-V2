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

DROP TABLE IF EXISTS SavedGamesPlayers;
DROP TABLE IF EXISTS SavedGamesTeams;
