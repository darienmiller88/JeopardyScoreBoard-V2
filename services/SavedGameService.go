package services

import (
	"context"

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
	return s.SavedGameRepository.AddSavedGame(ctx, savedGame)
}

func (s *SaveGameService) DeleteSavedGame(ctx context.Context, savedGameId bson.ObjectID) models.Result[*mongo.DeleteResult]{
	return s.SavedGameRepository.DeleteSavedGame(ctx, savedGameId)
}

//Validate each player in the array of players to gurauntee they all exist as real players in the database.
func (s *SaveGameService) validatePlayers(ctx context.Context, playersToValidate []models.Player) error{
	locationsResult := s.LocationRepository.GetAllLocations()

	if locationsResult.Err != nil {
		return locationsResult.Err
	}

	return  nil
}