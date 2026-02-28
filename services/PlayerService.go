package services

import (
	"net/http"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
	"JeopardyScoreBoardV2/utils"
)

type PlayerService interface {
	UpdatePlayerName(oldPlayerId string, firstName string, lastName string, locationName string) models.Result[models.Player]
	AddPlayerToLocation(locationName string, firstName string, lastName string) models.Result[models.Player]
	RemovePlayer(playerId string, locationName string) models.Result[models.Player]
	GetPlayersFromLocation(locationName string) models.Result[[]models.Player]
	GetAllPlayersFromAllLocations() models.Result[[]models.Player]
}

type PlayerServiceImpl struct {
	PlayerRepository repositories.PlayerRepository
}

// Update a players old name to be a new name.
func (p *PlayerServiceImpl) UpdatePlayerName(oldPlayerId string, firstName string, lastName string, locationName string) models.Result[models.Player] {
	player := models.Player{}

	//Set the player name to be the first name plus last name
	player.SetPlayerName(firstName, lastName)

	//Ensure the new name passes validation/is properly formatted
	if err := player.Validate(); err != nil {
		return utils.GetResult(err, http.StatusUnprocessableEntity, player)
	}

	//After confirming the following:
	//1. The new name is properly formatted
	//Update the players old name to be the new name.
	return p.PlayerRepository.UpdatePlayerName(oldPlayerId, player.PlayerName, locationName)
}

func (p *PlayerServiceImpl) AddPlayerToLocation(locationName string, firstName string, lastName string) models.Result[models.Player] {
	player := models.Player{}

	//set the player name by using the first name and last name.
	player.SetPlayerName(firstName, lastName)

	//Ensure the player name (first + last) passes all validation.
	if err := player.Validate(); err != nil {
		return utils.GetResult(err, http.StatusUnprocessableEntity, player)
	}

	//After confirming the following:
	//1. The new name is properly formatted (validated)
	//Add the new player to be database.
	return p.PlayerRepository.AddPlayerToLocation(locationName, player)
}

// Remove a player belonging to a certain location by using their name and location.
func (p *PlayerServiceImpl) RemovePlayer(playerId string, locationName string) models.Result[models.Player] {
	return p.PlayerRepository.RemovePlayer(playerId, locationName)
}

func (p *PlayerServiceImpl) GetAllPlayersFromAllLocations() models.Result[[]models.Player] {
	return p.PlayerRepository.GetAllPlayersFromAllLocations()
}

func (p *PlayerServiceImpl) GetPlayersFromLocation(locationName string) models.Result[[]models.Player] {
	return p.PlayerRepository.GetPlayersFromLocation(locationName)
}
