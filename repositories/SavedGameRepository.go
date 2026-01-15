package repositories

import (
	"JeopardyScoreBoardV2/models"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type SavedGameRepository interface {
	GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame]
	AddSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame]
	DeleteSavedGame(savedGameId int) models.Result[string]
	GetAllSavedGames() models.Result[[]models.SavedGame]
}

type sqlSavedGameRepository struct {
	db *sqlx.DB
}

// Receive new Instance of MongoPlayerCardRepository.
func GetSqlSavedGameRepository(newDB *sqlx.DB) *sqlSavedGameRepository {
	return &sqlSavedGameRepository{db: newDB}
}

// Get all Saved games from database.
func (s *sqlSavedGameRepository) GetAllSavedGames() models.Result[[]models.SavedGame] {
	return getResult(nil, http.StatusOK, []models.SavedGame{})
}

// Get all saved games played at a specific location.
func (s *sqlSavedGameRepository) GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame] {
	return getResult(nil, http.StatusOK, []models.SavedGame{})
}

// Delete a saved game
func (m *sqlSavedGameRepository) DeleteSavedGame(savedGameId int) models.Result[string] {
	return getResult(nil, http.StatusOK, "")
}

// Add a new saved game
func (m *sqlSavedGameRepository) AddSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame] {
	return getResult(nil, http.StatusOK, savedGame)
}
