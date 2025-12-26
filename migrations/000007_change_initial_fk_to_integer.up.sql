-- drop auto-incremnting behaviors from both FKs from the players table
ALTER TABLE Players ALTER COLUMN location_id DROP DEFAULT;
ALTER TABLE Players ALTER COLUMN team_id DROP DEFAULT;

-- Drop auto-incrementing from FK from SavedGames
ALTER TABLE SavedGames ALTER COLUMN location_id DROP DEFAULT;

-- Drop auto-incrementing from FK from Teams
ALTER TABLE Teams ALTER COLUMN location_id DROP DEFAULT;