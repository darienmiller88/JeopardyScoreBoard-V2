package repositories

import (
	"JeopardyScoreBoardV2/models"
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

const(
	push string = "$push"
	pull string = "$pull"
)

//LocationRepository interface to allow mocking when testing the service. The test can provide the service 
//a dummy implementation 
type LocationRepository interface{
	AddLocation(ctx context.Context)  models.Result[models.Location]
	GetLocations(ctx context.Context) models.Result[models.Location]
	GetLocation(ctx context.Context, locationName string) models.Result[models.Location]
}

type MongoLocationRepository struct{
	locationCollection *mongo.Collection
}

//Receive a new instance of Location repository using a mongo collection as the database. 
func GetMongoLocationCollection(newCollection *mongo.Collection) *MongoLocationRepository{
	return &MongoLocationRepository{ locationCollection: newCollection.Clone() }
}

func (m *MongoLocationRepository) AddLocation() models.Result[models.Location]{
	return models.Result[models.Location]{}
}

func (m *MongoLocationRepository) GetLocations() models.Result[models.Location]{
	return models.Result[models.Location]{}
}

func (m *MongoLocationRepository) GetLocation(locationName string) models.Result[models.Location]{
	return models.Result[models.Location]{}
}
