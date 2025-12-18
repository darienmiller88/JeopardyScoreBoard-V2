package services

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"

	"github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PlayerCardService struct{
	PlayerCardRepository repositories.PlayerCardRepository
	LocationRepository repositories.LocationRepository
}

func (p *PlayerCardService) UpdatePlayerName(ctx context.Context, locationName string, oldPlayerName string, newPlayerName string)models.Result[*mongo.UpdateResult]{
	playerNames := struct{
		NewPlayerName string 
		OldPlayerName string		
	}{ OldPlayerName: oldPlayerName, NewPlayerName: newPlayerName }

	//Validate the above struct to ensure both fields are included, and that the new name is between 3 and 30 chars.
	err := validation.ValidateStruct(&playerNames,
		validation.Field(&playerNames.NewPlayerName, validation.Required, validation.Length(3, 31)),
		validation.Field(&playerNames.OldPlayerName, validation.Required),
	)

	//If the first round of validation does not pass, return the following error.
	if err != nil{
		return models.Result[*mongo.UpdateResult]{ Err: err, StatusCode: http.StatusBadRequest }
	}
	
	//Afterwards, get the location to ensure the client sent one that actually exists. 
	locationResult := p.LocationRepository.GetLocation(ctx, locationName)

	// If not, return a 404 or 500.
	if locationResult.Err != nil {
		return models.Result[*mongo.UpdateResult]{ Err: locationResult.Err,StatusCode: http.StatusNotFound }
	}
	
	//Check to see if the new name the client is trying to send is already taken.
	if slices.Contains(locationResult.ResultData.Players, models.PlayerCard{ Name: playerNames.NewPlayerName }) {
		err := fmt.Errorf("player name '%s' is taken! Pick another one", playerNames.OldPlayerName)

		return models.Result[*mongo.UpdateResult]{ Err: err, StatusCode: http.StatusConflict }
	}
	
	//Finally, after checking to see if the names aren't blank, the new name is 3 < x < 31, and the new name
	//isn't taken, change the old name to the new name.
	return p.PlayerCardRepository.UpdatePlayerName(ctx, locationName, oldPlayerName, newPlayerName)
}

func (p *PlayerCardService) AddPlayerToLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult] {
	return p.PlayerCardRepository.AddPlayerToLocation(ctx, locationName, playerName)
}

func (p *PlayerCardService) RemovePlayerFromLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]{
	return  p.PlayerCardRepository.RemovePlayerFromLocation(ctx, locationName, playerName)
}