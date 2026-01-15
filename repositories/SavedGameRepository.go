package repositories

import (
	"JeopardyScoreBoardV2/models"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type SavedGameRepository interface {
	GetAllSavedGamesFromLocationDB(locationName string) models.Result[[]models.SavedGame]
	AddSavedGameDB(savedGame models.SavedGame)          models.Result[models.SavedGame]
	DeleteSavedGameDB(savedGameId int)                  models.Result[string]
	GetAllSavedGamesDB()                                models.Result[[]models.SavedGame]
}

type sqlSavedGameRepository struct {
	db *sqlx.DB
}

// Receive new Instance of MongoPlayerCardRepository.
func GetSqlSavedGameRepository(newDB *sqlx.DB) *sqlSavedGameRepository {
	return &sqlSavedGameRepository{db: newDB}
}

// Get all Saved games from database.
func (s *sqlSavedGameRepository) GetAllSavedGamesDB() models.Result[[]models.SavedGame] {
	return getResult(nil, http.StatusOK, []models.SavedGame{})
}

// Get all saved games played at a specific location.
func (s *sqlSavedGameRepository) GetAllSavedGamesFromLocationDB(locationName string) models.Result[[]models.SavedGame] {
	return getResult(nil, http.StatusOK, []models.SavedGame{})
}

// Delete a saved game
func (m *sqlSavedGameRepository) DeleteSavedGameDB(savedGameId int) models.Result[string] {
	return getResult(nil, http.StatusOK, "")
}

// Add a new saved game
func (m *sqlSavedGameRepository) AddSavedGameDB(savedGame models.SavedGame) models.Result[models.SavedGame] {
	return getResult(nil, http.StatusOK, savedGame)
}
