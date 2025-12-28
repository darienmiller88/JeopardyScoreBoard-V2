package repositories

import (
	"context"
	"net/http"
	"JeopardyScoreBoardV2/models"

	"github.com/jmoiron/sqlx"
)

const(
	push string = "$push"
	pull string = "$pull"
)

type PlayerRepository interface {
	UpdatePlayerName(locationName string, oldPlayerName string, newPlayerName string) models.Result[string]
	RemovePlayerFromLocation(locationName string, playerName string)                  models.Result[string]
	AddPlayerToLocation(locationName string, playerName string)                       models.Result[string]
	GetPlayersFromLocation(locationName string)                                       models.Result[[]models.Player]
	GetAllPlayersFromAllLocations()                                                   models.Result[[]models.Player]
}

type sqlPlayerRepository struct{
	db *sqlx.DB
}

//Receive new Instance of MongoPlayerCardRepository.
func GetNewMongoPlayerCardRepository(newDB *sqlx.DB) *sqlPlayerRepository{
	return &sqlPlayerRepository{ db: newDB }
}

//Function to update a players name for a given location.
func (s *sqlPlayerRepository) UpdatePlayerName(locationName string, oldPlayerName string, newPlayerName string) models.Result[models.Player]{
	return getResult(nil, 200, models.Player{})
}

//Add a single player to a given location.
func (s *sqlPlayerRepository) AddPlayerToLocation(locationName string, playerName string) models.Result[models.Player]{
	return getResult(nil, 200, models.Player{})
}

//Remove a single player from a given location.
func (s *sqlPlayerRepository) RemovePlayerFromLocation(ctx context.Context, locationName string, playerName string) models.Result[models.Player]{
	return getResult(nil, http.StatusOK, models.Player{})
}