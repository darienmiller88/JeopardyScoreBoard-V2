package repositories

import (
	"JeopardyScoreBoardV2/models"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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

//Get all Saved games from database.
func (m *MongoSavedGameRepository) GetAllSavedGames(ctx context.Context) models.Result[[]models.SavedGame]{
	findResult, err := m.savedGameCollection.Find(ctx, bson.D{})
	
	//If there was an error retrieving all documents, return the error, and a 500
	if err != nil {
		return models.Result[[]models.SavedGame]{ Err: err, StatusCode: http.StatusInternalServerError }
	}

	savedGames := []models.SavedGame{}

	//Unmarshall the mongo cursor into the array of saved games.
	if err := findResult.All(ctx, &savedGames); err != nil{
		return models.Result[[]models.SavedGame]{
			Err: err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//Return the retrieved saved games.
	return models.Result[[]models.SavedGame]{ ResultData: savedGames, StatusCode: http.StatusOK }
}

//Get all saved games played at a specific location.
func (m *MongoSavedGameRepository) GetAllSavedGamesFromLocation(ctx context.Context, locationName string) models.Result[[]models.SavedGame]{
	return models.Result[[]models.SavedGame]{}
}

//Delete a saved game
func (m *MongoSavedGameRepository) DeleteSavedGame(ctx context.Context, savedGameId primitive.ObjectID) models.Result[*mongo.DeleteResult]{
	return models.Result[*mongo.DeleteResult]{}
}

//Add a new saved game
func (m *MongoSavedGameRepository) AddSavedGame(ctx context.Context, savedGame models.SavedGame) models.Result[models.SavedGame]{
	return models.Result[models.SavedGame]{}
}