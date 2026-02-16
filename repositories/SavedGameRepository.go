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
	GetAllSavedGamesFromLocationDB(locationName string) models.Result[[]models.SavedGame]
	AddSavedGameDB(savedGame models.SavedGame) models.Result[models.SavedGame]
	DeleteSavedGameDB(savedGameId string) models.Result[string]
	GetAllSavedGamesDB() models.Result[[]models.SavedGame]
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

// Get all Saved games from database.
func (s *sqlSavedGameRepository) GetAllSavedGamesDB() models.Result[[]models.SavedGame] {
	savedGames := []models.SavedGame{}

	if err := s.db.Select(&savedGames, constants.GetAllSavedGames); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, []models.SavedGame{})
	}

	result := s.getSavedGamesWithDecryptedWinningPlayerName(savedGames)

	if result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, []models.SavedGame{})
	}

	return utils.GetResult(nil, http.StatusOK, savedGames)
}

//Get a saved game with a particular id
func (s *sqlSavedGameRepository) GetSavedGameById(savedGameId int) models.Result[models.SavedGame]{
	savedGame := models.SavedGame{}

	if err := s.db.Get(&savedGame, constants.GetSavedGameById, savedGame); err != nil {
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
func (s *sqlSavedGameRepository) DeleteSavedGameDB(savedGameId string) models.Result[string] {
	result, err := s.db.Exec(constants.DeleteSavedGame, savedGameId)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, "")
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, "")
	}

	if rowsAffected == 0 {
		return utils.GetResult(fmt.Errorf("saved game not found by id: %s", savedGameId), http.StatusNotFound, "")
	}

	return utils.GetResult(nil, http.StatusOK, fmt.Sprintf("Saved game %s deleted successfully", savedGameId))
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
