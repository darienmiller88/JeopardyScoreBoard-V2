package services

import (
	"net/http"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type PlayerService interface{
	UpdatePlayerName(locationName string, oldPlayerName string, newPlayerName string) models.Result[string]
	RemovePlayerFromLocation(locationName string, playerName string)                  models.Result[string]
	AddPlayerToLocation(locationName string, playerName string)                       models.Result[string]
	GetPlayersFromLocation(locationName string)                                       models.Result[[]models.Player]
	GetAllPlayersFromAllLocations()                                                   models.Result[[]models.Player]
}

type PlayerServiceImpl struct{
	PlayerRepository repositories.PlayerRepository
	LocationRepository repositories.LocationRepository
}

func (p *PlayerServiceImpl) UpdatePlayerName(locationName string, oldPlayerName string, newPlayerName string) models.Result[string]{
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
		return models.Result[string]{ Err: err, StatusCode: http.StatusBadRequest }
	}
	
	//Finally, after checking to see if the names aren't blank, the new name is 3 < x < 31, and the new name
	//isn't taken, change the old name to the new name.
	return p.PlayerRepository.UpdatePlayerName(locationName, oldPlayerName, newPlayerName)
}

func (p *PlayerServiceImpl) AddPlayerToLocation(locationName string, playerName string) models.Result[string] {
	
	
	return p.PlayerRepository.AddPlayerToLocation(locationName, playerName)
}

func (p *PlayerServiceImpl) RemovePlayerFromLocation(locationName string, playerName string) models.Result[string]{
	return  p.PlayerRepository.RemovePlayerFromLocation(locationName, playerName)
}