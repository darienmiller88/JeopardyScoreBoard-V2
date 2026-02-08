-- ==============================================
-- PLAYERS
-- ==============================================
ALTER TABLE players
DROP CONSTRAINT IF EXISTS uq_players_name_hash_location;



-- ==============================================
-- SAVEDGAMESPLAYERS
-- ==============================================
ALTER TABLE savedgamesplayers
DROP CONSTRAINT IF EXISTS uq_sgp_name_hash_game;



-- ==============================================
-- SAVEDGAMES
-- ==============================================
ALTER TABLE savedgames
DROP CONSTRAINT IF EXISTS uq_savedgames_winner_hash;
