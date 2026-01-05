package services

import (
	"net/http"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
)

type PlayerService interface{
	UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]
	RemovePlayerFromLocation(locationName string, playerName string)                  models.Result[models.Player]
	AddPlayerToLocation(locationName string, playerName string)                       models.Result[models.Player]
	GetPlayersFromLocation(locationName string)                                       models.Result[[]models.Player]
	GetAllPlayersFromAllLocations()                                                   models.Result[[]models.Player]
}

type PlayerServiceImpl struct{
	Repository repositories.PlayerRepository
}

func (p *PlayerServiceImpl) UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]{
	player := models.Player{ PlayerName: newPlayerName }

	if err := player.Validate(); err != nil {
		return models.Result[models.Player]{ Err: err, StatusCode: http.StatusUnprocessableEntity }
	}
	
	return p.Repository.UpdatePlayerName(oldPlayerName, newPlayerName)
}

func (p *PlayerServiceImpl) AddPlayerToLocation(locationName string, playerName string) models.Result[models.Player] {
	player := models.Player{ PlayerName: playerName }

	if err := player.Validate(); err != nil {
		return models.Result[models.Player]{ Err: err, StatusCode: http.StatusUnprocessableEntity }
	}
	
	return p.Repository.AddPlayerToLocation(locationName, player)
}

func (p *PlayerServiceImpl) RemovePlayerFromLocation(locationName string, playerName string) models.Result[models.Player]{
	return  p.Repository.RemovePlayerFromLocation(locationName, playerName)
}

func (p *PlayerServiceImpl) GetAllPlayersFromAllLocations() models.Result[[]models.Player]{
	return p.Repository.GetAllPlayersFromAllLocations()
}

func (p *PlayerServiceImpl) GetPlayersFromLocation(locationName string) models.Result[[]models.Player]{
	return p.Repository.GetPlayersFromLocation(locationName)
}