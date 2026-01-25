package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	push               string = "$push"
	pull               string = "$pull"
	uniqueKeyViolation string = "23505"
	notNullViolation   string = "23502"
)

type PlayerRepository interface {
	UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]
	AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player]
	GetPlayersFromLocation(locationName string) models.Result[[]models.Player]
	RemovePlayer(playerName string) models.Result[models.Player]
	GetAllPlayersFromAllLocations() models.Result[[]models.Player]

	ArePlayersValid(players []models.Player) models.Result[[]models.Player]	
}

type sqlPlayerRepository struct {
	db *sqlx.DB
}

// Receive new Instance of MongoPlayerCardRepository.
func GetSqlPlayerRepository(newDB *sqlx.DB) *sqlPlayerRepository {
	return &sqlPlayerRepository{db: newDB}
}

//determines if a list of players is valid (exists in the database)
func (s *sqlPlayerRepository) ArePlayersValid(players []models.Player) models.Result[[]models.Player]{
	validPlayers := []models.Player{}

	if err := s.db.Select(&validPlayers, constants.GetAllPlayers); err != nil {
		if err == sql.ErrNoRows {
			return utils.GetResult(fmt.Errorf("No players added yet"), http.StatusNotFound, validPlayers)
		}

		return utils.GetResult(err, http.StatusInternalServerError, validPlayers)
	}

	validPlayersMap := make(map[string]struct{}, len(validPlayers))

	for _, player := range validPlayers {
		validPlayersMap[player.PlayerName] = struct{}{}	
	}

	for _, player := range players {
		if _, ok := validPlayersMap[player.PlayerName]; !ok{
			return utils.GetResult(fmt.Errorf("player '%s' does not exist", player.PlayerName), http.StatusNotFound, []models.Player{})
		}
	}

	return utils.GetResult(nil, http.StatusOK, []models.Player{})
}

//Checks if a winning player exists
func (s *sqlSavedGameRepository) IsWinningPlayerValid(playerName string) models.Result[models.Player]{
	player := models.Player{}

	if err := s.db.Get(&player, constants.GetPlayerByName, playerName); err != nil{
		if err == sql.ErrNoRows {
			return utils.GetResult(fmt.Errorf("winning player '%s' does not exist", playerName), http.StatusNotFound, models.Player{})
		}

		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	return utils.GetResult(nil, http.StatusOK, models.Player{})
}

// Add a single player to a given location.
func (s *sqlPlayerRepository) AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player] {
	if err := s.db.Get(&player.ID, constants.InsertNewPlayerWithoutTeam, player.PlayerName, locationName); err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23502", "23503": // NOT NULL or FK violation
				return utils.GetResult(fmt.Errorf("no location '%s' found", locationName), http.StatusNotFound, models.Player{})
			case "23505": //unique key violation
				return utils.GetResult(fmt.Errorf("name %s is already taken", player.PlayerName), http.StatusConflict, models.Player{})
			}
		}

		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	return utils.GetResult(nil, http.StatusCreated, player)
}

// Function to update a players name for a given location.
func (s *sqlPlayerRepository) UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player] {
	result, err := s.db.Exec(constants.UpdatePlayerName, newPlayerName, oldPlayerName)

	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return utils.GetResult(fmt.Errorf("player name '%s' already exists", newPlayerName), http.StatusConflict, models.Player{})
		}

		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	numRowsAffected, _ := result.RowsAffected()

	if numRowsAffected == 0 {
		return utils.GetResult(fmt.Errorf("could not find player %s", oldPlayerName), http.StatusNotFound, models.Player{})
	}

	return utils.GetResult(nil, http.StatusOK, models.Player{PlayerName: newPlayerName})
}

// Remove a single player from a given location.
func (s *sqlPlayerRepository) RemovePlayer(playerName string) models.Result[models.Player] {
	result, err := s.db.Exec(constants.DeletePlayer, playerName)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	numRowsAffected, _ := result.RowsAffected()

	if numRowsAffected == 0 {
		return utils.GetResult(fmt.Errorf("could not find player %s", playerName), http.StatusNotFound, models.Player{})
	}

	return utils.GetResult(nil, http.StatusOK, models.Player{PlayerName: playerName})
}

func (s *sqlPlayerRepository) GetPlayersFromLocation(locationName string) models.Result[[]models.Player] {
	players := []models.Player{}

	if err := s.db.Select(&players, constants.GetAllPlayersFromLocation, locationName); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, players)
	}

	return utils.GetResult(nil, http.StatusOK, players)
}

func (s *sqlPlayerRepository) GetAllPlayersFromAllLocations() models.Result[[]models.Player] {
	players := []models.Player{}

	if err := s.db.Select(&players, constants.GetAllPlayers); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, players)
	}

	return utils.GetResult(nil, http.StatusOK, players)
}
