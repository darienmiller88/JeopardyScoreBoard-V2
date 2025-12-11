package repositories

import (
	"JeopardyScoreBoardV2/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SavedGameRepository interface {
	GetAllSavedGames(ctx context.Context) models.Result[[]models.SavedGame]
	GetAllSavedGamesFromLocation(ctx context.Context, locationName string) models.Result[[]models.SavedGame]
	DeleteSavedGame(ctx context.Context, savedGameId primitive.ObjectID) models.Result[*mongo.DeleteResult]
	AddSavedGame(ctx context.Context, savedGame models.SavedGame) models.Result[models.SavedGame]
}

type MongoSavedGameRepository struct{
	savedGameCollection *mongo.Collection
}

//Receive new Instance of MongoSavedGameRepository.
func GetNewMongoSavedGameRepository(newCollection *mongo.Collection) *MongoSavedGameRepository{
	return &MongoSavedGameRepository{ savedGameCollection: newCollection }
}

func (m *MongoSavedGameRepository) GetAllSavedGames(ctx context.Context) models.Result[[]models.SavedGame]{
	return models.Result[[]models.SavedGame]{}
}

func (m *MongoSavedGameRepository) GetAllSavedGamesFromLocation(ctx context.Context, locationName string) models.Result[[]models.SavedGame]{
	return models.Result[[]models.SavedGame]{}
}

func (m *MongoSavedGameRepository) DeleteSavedGame(ctx context.Context, savedGameId primitive.ObjectID) models.Result[*mongo.DeleteResult]{
	return models.Result[*mongo.DeleteResult]{}
}

func (m *MongoSavedGameRepository) AddSavedGame(ctx context.Context, savedGame models.SavedGame) models.Result[models.SavedGame]{
	return models.Result[models.SavedGame]{}
}