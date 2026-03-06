package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/encryption"
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/utils"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type SavedGameRepository interface {
	AddSavedGameDB(savedGame models.SavedGame) models.Result[models.SavedGame]
	DeleteSavedGameDB(savedGameId int) models.Result[string]
	
	//Get all saved games from a single location
	GetAllSavedGamesFromLocationDB(locationName string) models.Result[[]models.SavedGame]

	//get all saved games with all players who particiapted in each game
	GetAllSavedGamesDB() models.Result[[]models.SavedGame]

	//Get all players from a single saved game
	GetAllPlayersFromSavedGame(savedGameId string) models.Result[[]models.Player]

	//Get all saved games by id
	GetSavedGameById(savedGameId int) models.Result[models.SavedGame]
}

type sqlSavedGameRepository struct {
	db *sqlx.DB
	encryptionService *encryption.EncryptionService
}

// Receive new Instance of MongoPlayerCardRepository.
func GetSqlSavedGameRepository(newDB *sqlx.DB, 	encryptionService *encryption.EncryptionService) *sqlSavedGameRepository {
	return &sqlSavedGameRepository{db: newDB, encryptionService: encryptionService}
}

//Get all players from a single saved game
func (s *sqlSavedGameRepository) GetAllPlayersFromSavedGame(savedGameId string) models.Result[[]models.Player]{
	players := []models.Player{}

	if err := s.db.Select(&players, constants.GetAllPlayersFromSavedGame, savedGameId); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, []models.Player{})
	}

	return utils.GetResult(nil, http.StatusOK, players)
}

// Get all Saved games from database with players.
func (s *sqlSavedGameRepository) GetAllSavedGamesDB() models.Result[[]models.SavedGame] {
	savedGames := []models.SavedGame{}

	if err := s.db.Select(&savedGames, constants.GetAllSavedGames); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, []models.SavedGame{})
	}

	if len(savedGames) == 0 {
        return utils.GetResult(nil, http.StatusOK, savedGames)
    }

	result := s.getSavedGamesWithDecryptedWinningPlayerName(savedGames)

	if result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, []models.SavedGame{})
	}

	 // Collect IDs
    ids := []int{}
    for _, sg := range savedGames {
        ids = append(ids, sg.ID)
    }

    // Fetch players
    query, args, err := sqlx.In(`
        SELECT saved_game_id, player_id, player_score,
               player_name_encrypted, player_name_hash
        FROM savedgamesplayers
        WHERE saved_game_id IN (?)
    `, ids)

    if err != nil {
        return utils.GetResult(err, http.StatusInternalServerError, []models.SavedGame{})
    }

    query = s.db.Rebind(query)
    savedGamesPlayers := []models.SavedGamePlayer{}

    if err := s.db.Select(&savedGamesPlayers, query, args...); err != nil {
        return utils.GetResult(err, http.StatusInternalServerError, []models.SavedGame{})
    }

    playersByGame := make(map[int][]models.Player)

    for _, sgPlayer := range savedGamesPlayers {
        decrypted, _ := s.encryptionService.Decrypt(sgPlayer.NameEncrypted)
        player       := models.Player{
            ID:         sgPlayer.PlayerID,
            Score:      sgPlayer.PlayerScore,
            PlayerName: decrypted,
        }

        playersByGame[sgPlayer.SavedGameID] = append(playersByGame[sgPlayer.SavedGameID], player)
    }

    for i := range savedGames {
        savedGames[i].Players = playersByGame[savedGames[i].ID]

		//Set the flag to true for each game that is a player based game. If not, the flag is set to 
		//false for a team based game.
		if savedGames[i].WinningPlayerId.Valid {
			savedGames[i].IsPlayerGame = true
		}
    }

	return utils.GetResult(nil, http.StatusOK, savedGames)
}

//Get a saved game with a particular id
func (s *sqlSavedGameRepository) GetSavedGameById(savedGameId int) models.Result[models.SavedGame]{
	savedGame := models.SavedGame{}

	if err := s.db.Get(&savedGame, constants.GetSavedGameById, savedGameId); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	result := s.getSavedGamesWithDecryptedWinningPlayerName([]models.SavedGame{savedGame})

	if result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, models.SavedGame{})
	}

	return utils.GetResult(nil, http.StatusOK, result.ResultData[0])
}

// Get all saved games played at a specific location.
func (s *sqlSavedGameRepository) GetAllSavedGamesFromLocationDB(locationName string) models.Result[[]models.SavedGame] {
	savedGames := []models.SavedGame{}

	if err := s.db.Select(&savedGames, constants.GetAllSavedGamesFromLocation, locationName); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, []models.SavedGame{})
	}

	result := s.getSavedGamesWithDecryptedWinningPlayerName(savedGames)

	if result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, []models.SavedGame{})
	}

	return utils.GetResult(nil, http.StatusOK, savedGames)
}

// Delete a saved game
func (s *sqlSavedGameRepository) DeleteSavedGameDB(savedGameId int) models.Result[string] {
	result, err := s.db.Exec(constants.DeleteSavedGame, savedGameId)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, "")
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, "")
	}

	if rowsAffected == 0 {
		return utils.GetResult(fmt.Errorf("saved game not found by id: %d", savedGameId), http.StatusNotFound, "")
	}

	return utils.GetResult(nil, http.StatusOK, fmt.Sprintf("Saved game %d deleted successfully", savedGameId))
}

// Add a new saved game
func (s *sqlSavedGameRepository) AddSavedGameDB(savedGame models.SavedGame) models.Result[models.SavedGame] {
	if savedGame.IsPlayerGame {
		return s.addStandardSavedGame(savedGame)
	} else {
		return s.addTeamSavedGame(savedGame)
	}
}

func (s *sqlSavedGameRepository) addStandardSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame] {
	// Start transaction
	tx, err := s.db.Beginx() 

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	defer tx.Rollback()

	// Insert saved game
	err = tx.QueryRow(
		constants.InsertNewPlayerSavedGame,
		savedGame.TotalPoints,//$1
		savedGame.AveragePoints,//$2
		savedGame.WinningPlayerNameEncrypted, // $3 - encrypted winning player name
		savedGame.WinningPlayerNameHash,      // $4 - hashed winning player name
		savedGame.WinningPlayerId,//$5
		savedGame.LocationId,//$6
	).Scan(&savedGame.ID, &savedGame.CreatedAt, &savedGame.UpdatedAt)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	// Loop over and add each player to the junction table "savedgamesplayers"
	for _, player := range savedGame.Players {

		// Encrypt and hash each player's name
		encryptedPlayerName, err := s.encryptionService.Encrypt(player.PlayerName)

		if err != nil {
			return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
		}

		playerNameHash := encryption.NameHash(player.PlayerName)

		_, err = tx.Exec(
			constants.InsertPlayersForSavedGame,
			player.ID,            // $1
			savedGame.ID,         // $2
			player.Score,         // $3
			encryptedPlayerName,  // $4 - encrypted player name
			playerNameHash,       // $5 - hashed player name
		)

		if err != nil {
			return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	return utils.GetResult(nil, http.StatusCreated, savedGame)
}

func (s *sqlSavedGameRepository) addTeamSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame] {
	// Start transaction
	tx, err := s.db.Beginx()

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	defer tx.Rollback()

	// Insert saved game (team games don't have winning player names)
	err = tx.QueryRow(
		constants.InsertNewTeamSavedGame,
		savedGame.TotalPoints,
		savedGame.AveragePoints,
		savedGame.WinningTeamId,
		savedGame.LocationId,
	).Scan(&savedGame.ID)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	// Loop over and add each team to the junction table "savedgameteams"
	for _, team := range savedGame.Teams {
		_, err := tx.Exec(
			constants.InsertTeamsForSavedGame,
			team.ID,
			savedGame.ID,
			team.Score,
		)

		if err != nil {
			return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
		}
		
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	return utils.GetResult(nil, http.StatusCreated, savedGame)
}

func (s *sqlSavedGameRepository) getSavedGamesWithDecryptedWinningPlayerName(savedGames []models.SavedGame) models.Result[[]models.SavedGame]{
	//Decrypt the winning player name for each saved game 
	for i := range savedGames {
		decryptedName, err := s.encryptionService.Decrypt(savedGames[i].WinningPlayerNameEncrypted)

		if err != nil{
			return utils.GetResult(err, http.StatusInternalServerError, []models.SavedGame{})
		}

		savedGames[i].WinningPlayerName = decryptedName
	}

	return utils.GetResult(nil, http.StatusOK, savedGames)
}
