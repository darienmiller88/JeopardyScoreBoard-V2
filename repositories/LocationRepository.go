package repositories

import (
	"JeopardyScoreBoardV2/models"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const(
	push string = "$push"
	pull string = "$pull"
)

//LocationRepository interface to allow mocking when testing the service. The test can provide the service 
//a dummy implementation 
type LocationRepository interface{
	AddLocation(ctx context.Context, location models.Location) models.Result[models.Location]
	GetLocation(ctx context.Context, locationName string)      models.Result[models.Location]
	GetAllLocations(ctx context.Context)                       models.Result[[]models.Location]
}

type MongoLocationRepository struct{
	locationCollection *mongo.Collection
}

//Receive a new instance of Location repository using a mongo collection as the database. 
func GetMongoLocationCollection(newCollection *mongo.Collection) *MongoLocationRepository{
	return &MongoLocationRepository{ locationCollection: newCollection.Clone() }
}

//Add a new adapt location, which for now, I will not expose.
func (m *MongoLocationRepository) AddLocation(ctx context.Context, location models.Location) models.Result[models.Location]{
	insertOneResult, err := m.locationCollection.InsertOne(ctx, location)
	locationResult := models.Result[models.Location]{}

	//If there was an error inserting the location, return it along a 500
	if err != nil{
		locationResult.Err = err
		locationResult.StatusCode = http.StatusInternalServerError

		return locationResult
	}
	
	//Afterwards, take the inserted id of the location, add it to the payload, and return it with a 200
	location.ID = insertOneResult.InsertedID.(primitive.ObjectID)
	locationResult.ResultData = location
	locationResult.StatusCode = http.StatusOK

	return locationResult
}

//Retrieve all Locations from database
func (m *MongoLocationRepository) GetAllLocations(ctx context.Context) models.Result[[]models.Location]{
	findResult, err := m.locationCollection.Find(ctx, bson.D{})

	if err != nil {
		return models.Result[[]models.Location]{ StatusCode: http.StatusInternalServerError, Err: err }
	}

	locations := []models.Location{}

	if err := findResult.All(ctx, &locations); err != nil {
		return models.Result[[]models.Location]{ StatusCode: http.StatusInternalServerError, Err: err }
	}

	return models.Result[[]models.Location]{ StatusCode: http.StatusOK, ResultData: locations }
}

func (m *MongoLocationRepository) GetLocation(locationName string) models.Result[models.Location]{
	return models.Result[models.Location]{}
}
