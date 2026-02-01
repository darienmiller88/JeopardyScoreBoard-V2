package services

import (
	"fmt"
	"net/http"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
	"JeopardyScoreBoardV2/utils"
)

type PlayerService interface{
	UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]
	AddPlayerToLocation(locationName string, playerName string)  models.Result[models.Player]
	RemovePlayer(playerName string)                              models.Result[models.Player]
	GetPlayersFromLocation(locationName string)                  models.Result[[]models.Player]
	GetAllPlayersFromAllLocations()                              models.Result[[]models.Player]
}

type PlayerServiceImpl struct{
	Repository repositories.PlayerRepository
}

func (p *PlayerServiceImpl) UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]{
	player := models.Player{ PlayerName: newPlayerName }

	if err := player.Validate(); err != nil {
		return utils.GetResult(err, http.StatusUnprocessableEntity, player)
	}

	if oldPlayerName == newPlayerName {
		return utils.GetResult(fmt.Errorf("old and new names must be different"), http.StatusUnprocessableEntity, player)			
	}

	//ensure the new name isn't taken
	
	return p.Repository.UpdatePlayerName(oldPlayerName, newPlayerName)
}

func (p *PlayerServiceImpl) AddPlayerToLocation(locationName string, playerName string) models.Result[models.Player] {
	player := models.Player{ PlayerName: playerName }

	if err := player.Validate(); err != nil {
		return models.Result[models.Player]{ Err: err, StatusCode: http.StatusUnprocessableEntity }
	}
	
	//ensure the new name isn't taken

	return p.Repository.AddPlayerToLocation(locationName, player)
}

func (p *PlayerServiceImpl) RemovePlayer(playerName string) models.Result[models.Player]{
	return  p.Repository.RemovePlayer(playerName)
}

func (p *PlayerServiceImpl) GetAllPlayersFromAllLocations() models.Result[[]models.Player]{
	return p.Repository.GetAllPlayersFromAllLocations()
}

func (p *PlayerServiceImpl) GetPlayersFromLocation(locationName string) models.Result[[]models.Player]{
	return p.Repository.GetPlayersFromLocation(locationName)
}

// func (p *PlayerServiceImpl) isPlayerNameTaken(playerName string) models.Result[models.Player]{
// 	result := p.Repository.GetPlayerByName(playerName)

// 	//If the repo returned a row, it means a name is taken
// 	if result.ResultData.PlayerName != "" {
		
// 	}
// }