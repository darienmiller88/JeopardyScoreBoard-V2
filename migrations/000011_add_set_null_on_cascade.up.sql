ALTER TABLE savedgames DROP CONSTRAINT savedgames_player_id_fkey;
ALTER TABLE savedgames 
ADD CONSTRAINT savedgames_player_id_fkey FOREIGN KEY (player_id) REFERENCES Players(id) ON DELETE SET NULL;

ALTER TABLE savedgamesplayers DROP CONSTRAINT savedgamesplayers_player_id_fkey;
ALTER TABLE savedgamesplayers 
ADD CONSTRAINT savedgamesplayers_player_id_fkey FOREIGN KEY (player_id) REFERENCES Players(id) ON DELETE SET NULL;