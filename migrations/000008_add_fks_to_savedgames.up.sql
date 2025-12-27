-- Drop the following column as it is no longer needed.
ALTER TABLE savedgames DROP COLUMN is_team_game;

-- Add a column to store the name of the winner
ALTER TABLE savedgames 
ADD COLUMN winner_player_name VARCHAR(60);

-- Add winning team as a new column
ALTER TABLE savedgames ADD COLUMN winning_team_id INTEGER NOT NULL;

-- Add the column as an FK to teams
ALTER TABLE savedgames ADD CONSTRAINT savedgames_winning_team_id_fkey 
FOREIGN KEY (winning_team_id) REFERENCES Teams(id);