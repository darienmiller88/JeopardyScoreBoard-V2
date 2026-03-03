ALTER TABLE savedgames 
ADD CONSTRAINT uq_savedgames_winner_hash UNIQUE (winning_player_name_hash);