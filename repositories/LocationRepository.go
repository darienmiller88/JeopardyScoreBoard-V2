package repositories

import (
	"JeopardyScoreBoardV2/database"
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
	GetLocation(ctx context.Context, locationName string) models.Result[models.Location]
	GetAllLocations(ctx context.Context)                  models.Result[[]models.Location]
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
	findResult, err := database.GetLocationsCollection().Find(ctx, bson.D{})

	if err != nil {
		return models.Result[[]models.Location]{ StatusCode: http.StatusInternalServerError, Err: err }
	}

	locations := []models.Location{}

	if err := findResult.All(ctx, &locations); err != nil {
		return models.Result[[]models.Location]{ StatusCode: http.StatusInternalServerError, Err: err }
	}

	return models.Result[[]models.Location]{ StatusCode: http.StatusOK, ResultData: locations }
}

//Get one location from the database
func (m *MongoLocationRepository) GetLocation(ctx context.Context, locationName string) models.Result[models.Location]{
	location := &models.Location{}
	result   := models.Result[models.Location]{}
	err      := m.locationCollection.FindOne(
		ctx, 
		bson.D{{Key: "location_name", Value: locationName}},
	).Decode(&location)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			result.Err = fmt.Errorf("location \"%s\" does not exist. Please try another one", locationName)
			result.StatusCode = http.StatusNotFound
		} else {
			result.Err = err
			result.StatusCode = http.StatusInternalServerError
		}

		return result
	}

	result.ResultData = *location
	result.StatusCode = http.StatusOK

	return result
}