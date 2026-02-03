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

	return utils.GetResult(nil, http.StatusOK, savedGames)
}

// Get all saved games played at a specific location.
func (s *sqlSavedGameRepository) GetAllSavedGamesFromLocationDB(locationName string) models.Result[[]models.SavedGame] {
	savedGames := []models.SavedGame{}

	if err := s.db.Select(&savedGames, constants.GetAllSavedGamesFromLocation, locationName); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, []models.SavedGame{})
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
		savedGame.TotalPoints,
		savedGame.AveragePoints,
		savedGame.WinningPlayerName,
		savedGame.WinningPlayerId,
		savedGame.LocationId,
	).Scan(&savedGame.ID)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	// Loop over and add each player to the junction table "savedgameplayers"
	for _, player := range savedGame.Players {
		_, err := tx.Exec(
			constants.InsertPlayersForSavedGame,
			player.ID,         // $1
			savedGame.ID,      // $2
			player.Score,      // $3
			player.PlayerName, // $4
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

	// Insert saved game
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
