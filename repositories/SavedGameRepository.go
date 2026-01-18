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
func (s *sqlSavedGameRepository) DeleteSavedGameDB(savedGameId string) models.Result[string] {
	result, err := s.db.Exec(constants.DeleteSavedGame, savedGameId)

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
func (s *sqlSavedGameRepository) AddSavedGameDB(savedGame models.SavedGame) models.Result[models.SavedGame] {
	if savedGame.WinningPlayerId.Valid {
		
	}else{

	}
	
	return getResult(nil, http.StatusOK, savedGame)
}

func (s *sqlSavedGameRepository) addStandardSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame] {
	// Start transaction
	tx, err := s.db.Beginx()

	if err != nil {
		return getResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	defer tx.Rollback()

	// Insert saved game
	err = tx.QueryRow(
		constants.InsertNewPlayerSavedGame,
		savedGame.TotalPoints,
		savedGame.AveragePoints,
		savedGame.WinningPlayerName,
		savedGame.WinningPlayerName, // Used to lookup player id
		savedGame.LocationId,
	).Scan(&savedGame.ID)

	if err != nil {
		var pqErr *pq.Error
		// FK violation returns a 23503
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return getResult(fmt.Errorf("no location or winning player name found"), http.StatusNotFound, models.SavedGame{})
		}

		return getResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	// Loop over and add each player to the junction table "savedgameplayers"
	for _, player := range savedGame.Players {
		_, err := tx.Exec(
			constants.InsertPlayersForSavedGame,
			player.ID,
			savedGame.ID,
			player.Score,
			player.ID, // Lookup for player name
		)

		if err != nil {
			var pqErr *pq.Error

			if errors.As(err, &pqErr) {
				switch pqErr.Code {
				case "23503": // fk violation
					return getResult(fmt.Errorf("player with id %d not found", player.ID), http.StatusNotFound, models.SavedGame{})
				case "23505": // unique key violation
					return getResult(fmt.Errorf("duplicate player-id:game-id entry for player %d", player.ID), http.StatusConflict, models.SavedGame{})
				}
			}

			return getResult(err, http.StatusInternalServerError, models.SavedGame{})
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return getResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	return getResult(nil, http.StatusCreated, savedGame)
}

func (s *sqlSavedGameRepository) addTeamSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame] {
	// Start transaction
	tx, err := s.db.Beginx()

	if err != nil {
		return getResult(err, http.StatusInternalServerError, models.SavedGame{})
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
		var pqErr *pq.Error
		
		// FK violation returns a 23503
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return getResult(fmt.Errorf("no location or winning team found"), http.StatusNotFound, models.SavedGame{})
		}

		return getResult(err, http.StatusInternalServerError, models.SavedGame{})
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
			var pqErr *pq.Error

			if errors.As(err, &pqErr) {
				switch pqErr.Code {
				case "23503": // fk violation
					return getResult(fmt.Errorf("team with id %d not found", team.ID), http.StatusNotFound, models.SavedGame{})
				case "23505": // unique key violation
					return getResult(fmt.Errorf("duplicate team-id:game-id entry for team %d", team.ID), http.StatusConflict, models.SavedGame{})
				}
			}

			return getResult(err, http.StatusInternalServerError, models.SavedGame{})
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return getResult(err, http.StatusInternalServerError, models.SavedGame{})
	}

	return getResult(nil, http.StatusCreated, savedGame)
}