package services

import (
	"net/http"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type PlayerService interface{
	UpdatePlayerName(locationName string, oldPlayerName string, newPlayerName string) models.Result[models.Player]
	RemovePlayerFromLocation(locationName string, playerName string)                  models.Result[models.Player]
	AddPlayerToLocation(locationName string, playerName string)                       models.Result[models.Player]
	GetPlayersFromLocation(locationName string)                                       models.Result[[]models.Player]
	GetAllPlayersFromAllLocations()                                                   models.Result[[]models.Player]
}

type PlayerServiceImpl struct{
	Repository repositories.PlayerRepository
}

func (p *PlayerServiceImpl) UpdatePlayerName(locationName string, oldPlayerName string, newPlayerName string) models.Result[models.Player]{
	playerNames := struct{
		NewPlayerName string 
		OldPlayerName string		
	}{ OldPlayerName: oldPlayerName, NewPlayerName: newPlayerName }

	//Validate the above struct to ensure both fields are included, and that the new name is between 5 and 40 chars.
	err := validation.ValidateStruct(&playerNames,
		validation.Field(&playerNames.NewPlayerName, validation.Required, validation.Length(5, 40)),
		validation.Field(&playerNames.OldPlayerName, validation.Required),
	)

	//If the first round of validation does not pass, return the following error.
	if err != nil{
		return models.Result[models.Player]{ Err: err, StatusCode: http.StatusBadRequest }
	}
	
	//Finally, after checking to see if the names aren't blank, the new name is 5 <= x <= 40, 
	//change the old name to the new name.
	return p.Repository.UpdatePlayerName(locationName, oldPlayerName, newPlayerName)
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