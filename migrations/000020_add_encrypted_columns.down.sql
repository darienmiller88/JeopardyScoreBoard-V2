-- ==============================================
-- Drop columns and index from savedgamesplayers
-- ==============================================
ALTER TABLE savedgamesplayers
DROP INDEX idx_savedgamesplayers_name_hash;

ALTER TABLE savedgamesplayers
DROP COLUMN player_name_hash;

ALTER TABLE savedgamesplayers
DROP COLUMN player_name_encrypted;



-- ==============================================
-- Drop columns and index from savedgames
-- ==============================================
ALTER TABLE savedgames
DROP INDEX idx_savedgames_winner_hash;

ALTER TABLE savedgames
DROP COLUMN winning_player_name_encrypted;

ALTER TABLE savedgames
DROP COLUMN winning_player_name_hash;



-- ==============================================
-- Drop columns and index from players
-- ==============================================
ALTER TABLE players
DROP INDEX idx_players_name_hash;

ALTER TABLE players
DROP COLUMN player_name_encrypted;

ALTER TABLE players
DROP COLUMN player_name_hash;
