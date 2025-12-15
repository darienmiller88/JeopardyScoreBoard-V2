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
	GetAllSavedGamesFromLocation(ctx context.Context, locationName string) models.Result[[]models.SavedGame]
	DeleteSavedGame(ctx context.Context, savedGameId primitive.ObjectID)   models.Result[*mongo.DeleteResult]
	AddSavedGame(ctx context.Context, savedGame models.SavedGame)          models.Result[models.SavedGame]
	GetAllSavedGames(ctx context.Context)                                  models.Result[[]models.SavedGame]
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
		return models.Result[[]models.SavedGame]{ Err: err, StatusCode: http.StatusInternalServerError }
	}

	//Return the retrieved saved games.
	return models.Result[[]models.SavedGame]{ ResultData: savedGames, StatusCode: http.StatusOK }
}

//Get all saved games played at a specific location.
func (m *MongoSavedGameRepository) GetAllSavedGamesFromLocation(ctx context.Context, locationName string) models.Result[[]models.SavedGame]{
	savedGames      := []models.SavedGame{}
	findResult, err := m.savedGameCollection.Find(ctx, bson.D{
		{ Key: "location_name", Value: locationName },
	})
	
	//If finding all saved games from a certain location fails, send back the proper error.
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Result[[]models.SavedGame]{ Err: err, StatusCode: http.StatusNotFound }
		} 

		return models.Result[[]models.SavedGame]{ Err: err, StatusCode: http.StatusInternalServerError }
	}

	//Try and marshall all of the saved game documents in the saved games array, and send an error if it fails.
	if err := findResult.All(ctx, &savedGames); err != nil{
		return models.Result[[]models.SavedGame]{ Err: err, StatusCode: http.StatusInternalServerError }
	}

	//After checking for all error possibilities after retrieving all saved games, send back the data and a 200
	return models.Result[[]models.SavedGame]{ ResultData: savedGames, StatusCode: http.StatusOK }	
}

//Delete a saved game
func (m *MongoSavedGameRepository) DeleteSavedGame(ctx context.Context, savedGameId primitive.ObjectID) models.Result[*mongo.DeleteResult]{
	deleteResult, err := m.savedGameCollection.DeleteOne(ctx, bson.M{
		"_id": savedGameId,
	})

	//Try deleting a saved game by id, and add an error if it fails.
	if err != nil {
		return models.Result[*mongo.DeleteResult]{ Err: err, StatusCode: http.StatusBadRequest }
	}

	//Return the delete result and a 200
	return models.Result[*mongo.DeleteResult]{ ResultData: deleteResult, StatusCode: http.StatusOK }
}

//Add a new saved game
func (m *MongoSavedGameRepository) AddSavedGame(ctx context.Context, savedGame models.SavedGame) models.Result[models.SavedGame]{
	insertResult, err := m.savedGameCollection.InsertOne(ctx, &savedGame)

	if err != nil{ 
		return models.Result[models.SavedGame]{ Err: err, StatusCode: http.StatusInternalServerError }
	}
	
	//Finally, attach the id of the newly created saved game.
	savedGame.ID = insertResult.InsertedID.(primitive.ObjectID)
	
	//And return the saved game with its id along with a 200.
	return models.Result[models.SavedGame]{ ResultData: savedGame, StatusCode: http.StatusOK }
}