package constants

const(
	//CREATE
	InsertNewPlayerWithTeam string = `
		INSERT INTO players (player_name, location_id, team_id)
		VALUES(:player_name, :location_id, :team_id) RETURNING id
	`

	InsertNewPlayerWithoutTeam string = `
		INSERT INTO players (player_name, location_id)
		VALUES($1,
			(SELECT id FROM locations WHERE location_name=$2)
		) RETURNING id
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
		SELECT * FROM players
	`
	
	GetPlayerById string = `
		SELECT * FROM players WHERE id=$1
	`

	GetAllPlayersFromLocation string = `
		SELECT *
		FROM players 
		JOIN locations 
		ON locations.id=players.location_id
		WHERE locations.location_name=$1
	`

	//UPDATE
	UpdatePlayerName string = `
		UPDATE players SET player_name=$1 WHERE player_name=$2
	`

	//DELETE
)