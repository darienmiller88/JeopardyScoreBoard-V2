ALTER TABLE savedgames ADD COLUMN player_id INTEGER;
ALTER TABLE savedgames ADD CONSTRAINT savedgames_player_id_fkey 
FOREIGN KEY (player_id) REFERENCES Players(id)