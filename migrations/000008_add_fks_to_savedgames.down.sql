ALTER TABLE savedgames DROP CONSTRAINT savedgames_winning_team_id_fkey;

ALTER TABLE savedgames DROP COLUMN winning_team_id;

ALTER TABLE savedgames DROP COLUMN winner_player_name;

ALTER TABLE savedgames ADD COLUMN is_team_game BOOLEAN DEFAULT FALSE;