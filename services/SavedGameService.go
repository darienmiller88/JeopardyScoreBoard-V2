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

func (s *SaveGameService) GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGamesFromLocation(locationName)
}

func (s *SaveGameService) GetAllSavedGames() models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGames()
}

func (s *SaveGameService) AddSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame]{
	return s.SavedGameRepository.AddSavedGame(savedGame)
}

func (s *SaveGameService) DeleteSavedGame(savedGameId bson.ObjectID) models.Result[*mongo.DeleteResult]{
	return s.SavedGameRepository.DeleteSavedGame(savedGameId)
}