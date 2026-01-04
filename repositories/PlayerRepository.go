package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

const(
	push string = "$push"
	pull string = "$pull"
)

type PlayerRepository interface {
	UpdatePlayerName(locationName string, oldPlayerName string, newPlayerName string) models.Result[models.Player]
	RemovePlayerFromLocation(locationName string, playerName string)                  models.Result[models.Player]
	AddPlayerToLocation(locationName string, playerName string)                       models.Result[models.Player]
	GetPlayersFromLocation(locationName string)                                       models.Result[[]models.Player]
	GetAllPlayersFromAllLocations()                                                   models.Result[[]models.Player]
}

type sqlPlayerRepository struct{
	db *sqlx.DB
}

//Receive new Instance of MongoPlayerCardRepository.
func GetSqlPlayerCardRepository(newDB *sqlx.DB) *sqlPlayerRepository{
	return &sqlPlayerRepository{ db: newDB }
}

//Add a single player to a given location.
func (s *sqlPlayerRepository) AddPlayerToLocation(locationName string, playerName string) models.Result[models.Player]{
	statement, err := s.db.PrepareNamed(constants.InsertNewPlayerWithoutTeam)

	if err != nil {
		return getResult(err, http.StatusInternalServerError, models.Player{})
	}

	player := models.Player{}

	if err := statement.Get(&player.ID, player); err != nil{
		if err == sql.ErrNoRows {
			return getResult(fmt.Errorf("No location %s found", locationName), http.StatusNotFound, models.Player{})
		}
		
		return getResult(err, http.StatusInternalServerError, models.Player{})	
	}
	
	return getResult(nil, http.StatusOK, player)
}

//Function to update a players name for a given location.
func (s *sqlPlayerRepository) UpdatePlayerName(locationName string, oldPlayerName string, newPlayerName string) models.Result[models.Player]{
	return getResult(nil, 200, models.Player{})
}

//Remove a single player from a given location.
func (s *sqlPlayerRepository) RemovePlayerFromLocation(locationName string, playerName string) models.Result[models.Player]{
	return getResult(nil, http.StatusOK, models.Player{})
}

func (s *sqlPlayerRepository) GetPlayersFromLocation(locationName string) models.Result[[]models.Player]{
	return getResult(nil, 200, []models.Player{})
}

func (s *sqlPlayerRepository) GetAllPlayersFromAllLocations() models.Result[[]models.Player]{
	return getResult(nil, 200, []models.Player{})
}