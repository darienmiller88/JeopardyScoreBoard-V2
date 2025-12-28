package constants

const(
	//CREATE
	InsertNewPlayer string = `
		INSERT INTO players (player_name, is_active, left_at, location_id, team_id)
		VALUES(:player_name, :is_active, :left_at, :location_id, :team_id) RETURNING id
	`

	InsertNewSavedGame string = `
		INSERT INTO SavedGames (winner, total_score, average_score, is_team_game)
	`

	//READ
	GetAllLocations string = `
		SELECT location_name FROM locations
	`

	GetLocation string = `
		SELECT location_name FROM locations WHERE location_name=$1
	`

	//UPDATE


	//DESTROY
)