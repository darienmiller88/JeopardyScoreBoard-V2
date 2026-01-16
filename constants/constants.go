package constants

const(
	//CREATE
	InsertNewPlayerWithTeam string = `
		INSERT INTO players (player_name, location_id, team_id)
		VALUES(:player_name, :location_id, :team_id) RETURNING id
	`

	InsertNewPlayerWithoutTeam string = `
		INSERT INTO players (player_name, location_id)
		VALUES(
			$1,
			(SELECT id FROM locations WHERE location_name=$2)
		) RETURNING id
	`

	InsertNewPlayerSavedGame string = `
		INSERT INTO SavedGames (total_score, average_score, winning_player_name, winning_player_id, location_id)
		VALUES(
			$1, 
			$2,
			$3,
			(SELECT id FROM players WHERE player_name=$4),
			(SELECT id FROM locations WHERE location_name=$5)
		) RETURNING id
	`

	InsertNewTeamSavedGame string = `
		INSERT INTO SavedGames (total_score, average_score, winning_team_id, location_id)
		VALUES(
			$1,
			$2,
			$3,
			(SELECT id FROM locations WHERE location_name=$4)
		) RETURNING id
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
		SELECT players.*
		FROM players 
		JOIN locations 
		ON locations.id=players.location_id
		WHERE locations.location_name=$1
	`
	
	//Get saved games
	GetAllSavedGames string = `
		SELECT * FROM saved
	`

	GetAllSavedGamesFromLocation string = `
	
	`

	//UPDATE
	UpdatePlayerName string = `
		UPDATE players SET player_name=$1 WHERE player_name=$2
	`

	//DELETE
	DeletePlayer string = `
		DELETE FROM players WHERE player_name=$1
	`

	DeleteSavedGame string = `
	
	`

)