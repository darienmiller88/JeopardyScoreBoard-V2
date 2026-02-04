-- ==============================================
-- Drop indexes first (standalone)
-- ==============================================
DROP INDEX IF EXISTS idx_savedgamesplayers_name_hash;
DROP INDEX IF EXISTS idx_savedgames_winner_hash;
DROP INDEX IF EXISTS idx_players_name_hash;

-- ==============================================
-- Drop columns from savedgamesplayers
-- ==============================================
ALTER TABLE savedgamesplayers
DROP COLUMN IF EXISTS player_name_hash,
DROP COLUMN IF EXISTS player_name_encrypted;

-- ==============================================
-- Drop columns from savedgames
-- ==============================================
ALTER TABLE savedgames
DROP COLUMN IF EXISTS winning_player_name_encrypted,
DROP COLUMN IF EXISTS winning_player_name_hash;

-- ==============================================
-- Drop columns from players
-- ==============================================
ALTER TABLE players
DROP COLUMN IF EXISTS player_name_encrypted,
DROP COLUMN IF EXISTS player_name_hash;
