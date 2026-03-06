package constants

const(
	//CREATE


	//////////////////////////
	//Insert Player Queries
	//////////////////////////
	InsertNewPlayerWithTeam string = `
		INSERT INTO players (player_name, location_id, team_id)
		VALUES($1, $2, $3) RETURNING id
	`

	InsertNewPlayerWithoutTeam string = `
		INSERT INTO players (player_name_encrypted, player_name_hash, location_id)
		VALUES(
			$1,
			$2,
			(SELECT id FROM locations WHERE location_name=$3)
		) RETURNING id, created_at, updated_at, location_id
	`


	//////////////////////////////
	//Insert saved game Queries
	/////////////////////////////

	InsertNewPlayerSavedGame string = `
		INSERT INTO SavedGames (
			total_score,
			average_score,
			winning_player_name_encrypted,
			winning_player_name_hash,
			winning_player_id,
			location_id
		)
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		) RETURNING id, created_at, updated_at
	`

	InsertNewTeamSavedGame string = `
		INSERT INTO SavedGames (total_score, average_score, winning_team_id, location_id)
		VALUES(
			$1,
			$2,
			$3,
			$4
		) RETURNING id
	`

	InsertPlayersForSavedGame string = `
		INSERT INTO SavedGamesPlayers (
			player_id,
			saved_game_id,
			player_score,
			player_name_encrypted,
			player_name_hash
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
	`

	InsertTeamsForSavedGame string = `
		INSERT INTO SavedGamesTeams (team_id, saved_game_id, team_score)
		VALUES(
			$1,
			$2,
			$3
		)
	`
	//////////////////////////
	//READ
	/////////////////////////


	//////////////////////
	//Get Locations Queries
	//////////////////////

	//Get location(s)
	GetAllLocations string = `
		SELECT location_name FROM locations
	`

	GetLocation string = `
		SELECT location_name FROM locations WHERE location_name=$1
	`

	GetLocationById string = `
		SELECT location_name FROM locations WHERE id=$1
	`

	//////////////////////
	//Get Players Queries
	//////////////////////

	//Get player(s)
	GetAllPlayers string = `
		SELECT * FROM players
	`
	
	GetPlayerById string = `
		SELECT * FROM players WHERE id=$1
	`

	GetPlayerByName string = `
		SELECT * FROM players WHERE player_name_hash=$1
	`

	GetPlayersByNames string = `
		SELECT * FROM players WHERE player_name = ANY($1)
	`

	GetAllPlayersFromLocation string = `
		SELECT players.*
		FROM players 
		JOIN locations 
		ON locations.id=players.location_id
		WHERE locations.location_name=$1
	`

	

	//////////////////////
	//Get Teams queries
	//////////////////////

	//Get all team names, which are just the location names
	GetAllTeamsByName string = `
		SELECT locations.location_name
		FROM locations
		JOIN teams
		ON locations.id=teams.location_id
	`

	//Get all players on a certain team
	GetAllPlayersOnTeam string = `
		SELECT * FROM players WHERE team_id = $1
	`

	GetTeamById string = `
		SELECT * FROM teams WHERE id=$1
	`

	//Get a certain team by an id
	GetAllTeams string = `
		SELECT * FROM teams
	`

	GetTeamsByIds string = `
		SELECT * FROM teams WHERE id = ANY($1)
	`

	//////////////////////////////
	// Get Saved games
	/////////////////////////////
	
	//Get all saved games with the location name and winning player score
	GetAllSavedGames string = `
		SELECT 
			sg.*,
			l.location_name,
			sgp.player_score AS winning_player_score
		FROM savedgames sg
		JOIN locations l 
			ON sg.location_id = l.id
		JOIN savedgamesplayers sgp
			ON sgp.saved_game_id = sg.id 
			AND sgp.player_id = sg.winning_player_id
		ORDER BY sg.created_at DESC;
	`

	//Get all saved games from a certain location
	GetAllSavedGamesFromLocation string = `
		SELECT * FROM savedgames WHERE location_id=(SELECT id FROM locations WHERE location_name=$1)
	`

	GetSavedGameById string = `
		SELECT * FROM savedgames WHERE id=$1
	`

	GetAllPlayersFromSavedGame string = `
		SELECT * FROM savedgamesplayers WHERE saved_game_id=$1
	`

	/////////////////////////////
	//UPDATE
	/////////////////////////////

	//Update a players name
	UpdatePlayerName string = `
		UPDATE players
		SET
			updated_at=NOW(),
			player_name_encrypted = $1,
			player_name_hash = $2
		WHERE
			id = $3
			AND location_id = (
				SELECT id FROM locations WHERE location_name = $4
			);
	`


	//==========================
	// DELETE 
	//==========================

	//DELETE
	DeletePlayer string = `
		DELETE FROM players WHERE id=$1 AND location_id=(SELECT id FROM locations WHERE location_name=$2)
	`

	DeleteSavedGame string = `
		DELETE FROM savedgames WHERE id=$1
	`

)