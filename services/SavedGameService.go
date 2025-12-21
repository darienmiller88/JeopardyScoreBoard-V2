package services

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
)

type SaveGameService struct{
	SavedGameRepository repositories.SavedGameRepository
	LocationRepository  repositories.LocationRepository
}

func (s *SaveGameService) GetAllSavedGamesFromLocation(ctx context.Context, locationName string) models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGamesFromLocation(ctx, locationName)
}

func (s *SaveGameService) GetAllSavedGames(ctx context.Context) models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGames(ctx)
}

func (s *SaveGameService) AddSavedGame(ctx context.Context, savedGame models.SavedGame) models.Result[models.SavedGame]{
	//First, initialize the saved game, filling its created at field.
	savedGame.InitCreatedAtAndUpdatedAt()

	//Calculate the total amounts
	savedGame.CalcTotalPoints()

	//
	savedGame.CalcAveragePoints()



	return s.SavedGameRepository.AddSavedGame(ctx, savedGame)
}

func (s *SaveGameService) DeleteSavedGame(ctx context.Context, savedGameId bson.ObjectID) models.Result[*mongo.DeleteResult]{
	return s.SavedGameRepository.DeleteSavedGame(ctx, savedGameId)
}

//Validate each player in the array of players to gurauntee they all exist as real players in the database.
func (s *SaveGameService) validatePlayers(ctx context.Context, playersToValidate []models.PlayerCard) error{
	locationsResult := s.LocationRepository.GetAllLocations(ctx)

	if locationsResult.Err != nil {
		return locationsResult.Err
	}

	//create an array of players to extract all of the real players, so we can compare the players the client sent.
	realPlayers := []models.PlayerCard{}

	//Extract all of the players from all of the locations
	for _, location := range locationsResult.ResultData {
		realPlayers = append(realPlayers, location.Players...)
	}

	uniquePlayerNames := make(map[string]struct{})

	//Add each player to a map for easier indexing when comparing the list players sent here by the client.
	for _, player := range realPlayers{
		uniquePlayerNames[player.Name] = struct{}{}
	}

	//Check to see if any player in the list of players sent by the client exists.
	for _, player := range playersToValidate{
		if _, exists := uniquePlayerNames[player.Name]; !exists {
			return fmt.Errorf("player '%s' does not exist", player.Name)
		}
	}
	
	return  nil
}