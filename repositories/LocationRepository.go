package repositories

import (
	"JeopardyScoreBoardV2/models"
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


//LocationRepository interface to allow mocking when testing the service. The test can provide the service
//a dummy implementation
type LocationRepository interface{
	GetPlayersFromLocation(ctx context.Context, locationName string) models.Result[[]models.PlayerCard]
	GetLocation(ctx context.Context, locationName string)            models.Result[models.Location]
	GetAllPlayersFromAllLocations(ctx context.Context)               models.Result[[]models.PlayerCard]
	GetAllLocations(ctx context.Context)                             models.Result[[]models.Location]
}

type MongoLocationRepository struct{
	locationCollection *mongo.Collection
}

//Receive a new instance of Location repository using a mongo collection as the database. 
func GetMongoLocationCollection(newCollection *mongo.Collection) *MongoLocationRepository{
	return &MongoLocationRepository{ locationCollection: newCollection.Clone() }
}

//Retrieve all Locations from database
func (m *MongoLocationRepository) GetAllLocations(ctx context.Context) models.Result[[]models.Location]{
	findResult, err := m.locationCollection.Find(ctx, bson.D{})

	if err != nil {
		return getResult(err, http.StatusInternalServerError, []models.Location{})
	}

	locations := []models.Location{}

	if err := findResult.All(ctx, &locations); err != nil {
		return getResult(err, http.StatusInternalServerError, []models.Location{})
	}

	return getResult(nil, http.StatusOK, locations)
}

//Get one location from the database
func (m *MongoLocationRepository) GetLocation(ctx context.Context, locationName string) models.Result[models.Location]{
	location := models.Location{}
	err      := m.locationCollection.FindOne(ctx, bson.D{{Key: "location_name", Value: locationName}}).Decode(&location)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			err := fmt.Errorf("location \"%s\" does not exist. Please try another one", locationName)
			
			return getResult(err, http.StatusNotFound, models.Location{})
		} 

		return getResult(err, http.StatusInternalServerError, models.Location{})
	}

	return getResult(nil, http.StatusOK, location)
}

//Return the array of players from a particular location
func (m *MongoLocationRepository) GetPlayersFromLocation(ctx context.Context, locationName string) models.Result[[]models.PlayerCard]{
	locationResult := m.GetLocation(ctx, locationName)

	//Return the error from the above call if an error occurs
	if locationResult.Err != nil {
		return getResult(locationResult.Err, http.StatusInternalServerError, []models.PlayerCard{})
	}

	//IF not, return all of the players for for that location.
	return getResult(nil, http.StatusOK, locationResult.ResultData.Players)
}

//Return every single player at every single location.
func (m *MongoLocationRepository) GetAllPlayersFromAllLocations(ctx context.Context) models.Result[[]models.PlayerCard]{
	locationsResult := m.GetAllLocations(ctx)	

	if locationsResult.Err != nil {
		return getResult(locationsResult.Err, locationsResult.StatusCode, []models.PlayerCard{})
	}

	players := []models.PlayerCard{}

	//Range over each location, and extract all of the players for each location.
	for _, location := range locationsResult.ResultData {
		players = append(players, location.Players...)
	}

	//Return the list of all players in the database.
	return getResult(nil, http.StatusOK, players)
}

//Helper function to allow repos to send result payloads with less text.
func getResult[T any](err error, statusCode int, payload T) models.Result[T] {
	return models.Result[T]{
		StatusCode: statusCode,
		Err: err,
		ResultData: payload,
	}
}