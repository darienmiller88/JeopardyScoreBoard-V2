package constants

const(
	//CREATE
	InsertNewPlayer string = `
		INSERT INTO players (player_name, location_id, team_id)
		VALUES(:player_name, :location_id, :team_id) RETURNING id
	`

	InsertNewSavedGame string = `
		INSERT INTO SavedGames (winner, total_score, average_score, is_team_game)
	`

	//////////////////////////
	//READ
	/////////////////////////

	//Get location(s)
	GetAllLocations string = `
		SELECT location_name FROM locations
	`

	GetLocation string = `
		SELECT location_name FROM locations WHERE location_name=$1
	`

	//Get player(s)
	GetAllPlayers string = `
		SELECT player_name FROM players
	`

	//UPDATE


	//DESTROY
)