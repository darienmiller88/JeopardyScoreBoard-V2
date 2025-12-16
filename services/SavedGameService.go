package services

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
)

type SaveGameService struct{
	Repository repositories.SavedGameRepository
}

func (s *SaveGameService) GetAllSavedGamesFromLocation(ctx context.Context, locationName string) models.Result[[]models.SavedGame]{
	return s.Repository.GetAllSavedGamesFromLocation(ctx, locationName)
}

func (s *SaveGameService) GetAllSavedGames(ctx context.Context) models.Result[[]models.SavedGame]{
	return s.Repository.GetAllSavedGames(ctx)
}

func (s *SaveGameService) AddSavedGame(ctx context.Context, savedGame models.SavedGame) models.Result[models.SavedGame]{
	return  s.Repository.AddSavedGame(ctx, savedGame)
}

func (s *SaveGameService) DeleteSavedGame(ctx context.Context, savedGameId bson.ObjectID) models.Result[*mongo.DeleteResult]{
	return s.Repository.DeleteSavedGame(ctx, savedGameId)
}

