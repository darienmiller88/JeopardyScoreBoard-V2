ALTER TABLE players
ADD COLUMN player_name_encrypted BYTEA,
ADD COLUMN player_name_hash BYTEA;

CREATE INDEX idx_players_name_hash ON players(player_name_hash);

ALTER TABLE savedgames
ADD COLUMN winning_player_name_encrypted BYTEA,
ADD COLUMN winning_player_name_hash BYTEA;

CREATE INDEX idx_savedgames_winner_hash ON savedgames(winning_player_name_hash);

ALTER TABLE savedgamesplayers
ADD COLUMN player_name_encrypted BYTEA,
ADD COLUMN player_name_hash BYTEA;

CREATE INDEX idx_savedgamesplayers_name_hash ON savedgamesplayers(player_name_hash);