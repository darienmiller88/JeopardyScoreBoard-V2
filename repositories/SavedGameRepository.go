package repositories

import (
	"JeopardyScoreBoardV2/models"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SavedGameRepository interface {
	GetAllSavedGamesFromLocation(ctx context.Context, locationName string) models.Result[[]models.SavedGame]
	DeleteSavedGame(ctx context.Context, savedGameId bson.ObjectID)        models.Result[*mongo.DeleteResult]
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
		return getResult(err, http.StatusInternalServerError, []models.SavedGame{})
	}

	savedGames := []models.SavedGame{}

	//Unmarshall the mongo cursor into the array of saved games.
	if err := findResult.All(ctx, &savedGames); err != nil{
		return getResult(err, http.StatusInternalServerError, []models.SavedGame{})
	}

	//Return the retrieved saved games.
	return getResult(nil, http.StatusOK, savedGames)
}

//Get all saved games played at a specific location.
func (m *MongoSavedGameRepository) GetAllSavedGamesFromLocation(ctx context.Context, locationName string) models.Result[[]models.SavedGame]{
	savedGames      := []models.SavedGame{}
	findResult, err := m.savedGameCollection.Find(ctx, bson.D{{ Key: "location_name", Value: locationName }})
	
	//If finding all saved games from a certain location fails, send back the proper error.
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return getResult(err, http.StatusNotFound, []models.SavedGame{})
		} 

		return getResult(err, http.StatusInternalServerError, []models.SavedGame{})
	}

	//Try and marshall all of the saved game documents in the saved games array, and send an error if it fails.
	if err := findResult.All(ctx, &savedGames); err != nil{
		return getResult(err, http.StatusInternalServerError, []models.SavedGame{})
	}

	//After checking for all error possibilities after retrieving all saved games, send back the data and a 200
	return getResult(nil, http.StatusOK, savedGames)
}

//Delete a saved game
func (m *MongoSavedGameRepository) DeleteSavedGame(ctx context.Context, savedGameId bson.ObjectID) models.Result[*mongo.DeleteResult]{
	deleteResult, err := m.savedGameCollection.DeleteOne(ctx, bson.M{
		"_id": savedGameId,
	})

	//Try deleting a saved game by id, and add an error if it fails.
	if err != nil {
		return getResult(err, http.StatusInternalServerError, &mongo.DeleteResult{})
	}

	//Return the delete result and a 200
	return getResult(nil, http.StatusOK, deleteResult)
}

//Add a new saved game
func (m *MongoSavedGameRepository) AddSavedGame(ctx context.Context, savedGame models.SavedGame) models.Result[models.SavedGame]{
	
	
	//And return the saved game with its id along with a 200.
	return getResult(nil, http.StatusOK, savedGame)
}