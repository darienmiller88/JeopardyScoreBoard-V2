-- ==============================================
-- PLAYERS
-- Make (player_name_hash, location_id) unique
-- ==============================================

-- Drop old index if it exists
DROP INDEX IF EXISTS idx_players_name_hash;

-- Add proper unique constraint
ALTER TABLE players
ADD CONSTRAINT uq_players_name_hash_location
UNIQUE (player_name_hash, location_id);



-- ==============================================
-- SAVEDGAMESPLAYERS
-- Make (player_name_hash, saved_game_id) unique
-- ==============================================

DROP INDEX IF EXISTS idx_savedgamesplayers_name_hash;

ALTER TABLE savedgamesplayers
ADD CONSTRAINT uq_sgp_name_hash_game
UNIQUE (player_name_hash, saved_game_id);



-- ==============================================
-- SAVEDGAMES
-- Make winning_player_name_hash unique
-- ==============================================

DROP INDEX IF EXISTS idx_savedgames_winner_hash;

ALTER TABLE savedgames
ADD CONSTRAINT uq_savedgames_winner_hash
UNIQUE (winning_player_name_hash);