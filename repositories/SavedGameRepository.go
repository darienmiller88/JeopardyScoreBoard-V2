package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SavedGameRepository interface {
	GetAllSavedGamesFromLocationDB(locationName string) models.Result[[]models.SavedGame]
	AddSavedGameDB(savedGame models.SavedGame)          models.Result[models.SavedGame]
	DeleteSavedGameDB(savedGameId string)               models.Result[string]
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
	var savedGames []models.SavedGame
	
	if err := s.db.Select(&savedGames, constants.GetAllSavedGames); err != nil {
		return getResult(err, http.StatusInternalServerError, []models.SavedGame{})
	}
	
	return getResult(nil, http.StatusOK, savedGames)
}

// Get all saved games played at a specific location.
func (s *sqlSavedGameRepository) GetAllSavedGamesFromLocationDB(locationName string) models.Result[[]models.SavedGame] {
	var savedGames []models.SavedGame
	
	if err := s.db.Select(&savedGames, constants.GetAllSavedGamesFromLocation, locationName); err != nil {
		return getResult(err, http.StatusInternalServerError, []models.SavedGame{})
	}
	
	return getResult(nil, http.StatusOK, savedGames)
}

// Delete a saved game
func (m *sqlSavedGameRepository) DeleteSavedGameDB(savedGameId string) models.Result[string] {
	result, err := m.db.Exec(constants.DeleteSavedGame, savedGameId)

	if err != nil {
		return getResult(err, http.StatusInternalServerError, "")
	}
	
	rowsAffected, err := result.RowsAffected()
	
	if err != nil {
		return getResult(err, http.StatusInternalServerError, "")
	}
	
	if rowsAffected == 0 {
		return getResult(fmt.Errorf("saved game not found by id: %s", savedGameId), http.StatusNotFound, "")
	}
	
	return getResult(nil, http.StatusOK, fmt.Sprintf("Saved game %s deleted successfully", savedGameId))
}

// Add a new saved game
func (m *sqlSavedGameRepository) AddSavedGameDB(savedGame models.SavedGame) models.Result[models.SavedGame] {
	if savedGame.WinningPlayerId.Valid {
		err := m.db.Get(
			&savedGame.ID, //Retrieve the id of the newly created saved game
			constants.InsertNewPlayerSavedGame, 
			savedGame.TotalPoints,
			savedGame.AveragePoints,
			savedGame.WinningPlayerName, 
			savedGame.WinningPlayerName,//Used to lookup player id
			savedGame.LocationId,
		)

		if err != nil {
			var pqErr *pq.Error

			//FK violation returns a 23503
			if errors.As(err, &pqErr) && pqErr.Code == "23503" {
				return getResult(fmt.Errorf("no location or "), http.StatusNotFound, models.SavedGame{})
			}

			return getResult(err, http.StatusInternalServerError, models.SavedGame{})
		} 

		for _, player := range savedGame.Players {
			var pqErr *pq.Error
			_, err := m.db.Exec(
				constants.InsertPlayersForSavedGame,
				player.PlayerName,
				savedGame.ID,
				player.Score,
				
			)

			//FK violation returns a 23503
			if errors.As(err, &pqErr) && pqErr.Code == "23503" {
				return getResult(fmt.Errorf("no player with id %s found", ), http.StatusNotFound, models.SavedGame{})
			}
		}
	}else{

	}
	
	return getResult(nil, http.StatusOK, savedGame)
}
