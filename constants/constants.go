package constants

const(
	//CREATE
	insertNewPlayer string = `
		INSERT INTO players (player_name, is_active, left_at, location_id, team_id)
		VALUES(:player_name, :is_active, :left_at, :location_id, :team_id) RETURNING id
	`

	insertNewSavedGame string = `
		INSERT INTO SavedGames (winner, total_score, average_score, is_team_game)
	`

	//READ


	//UPDATE


	//DESTROY
)