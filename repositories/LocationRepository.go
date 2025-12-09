package repositories

import(
	"JeopardyScoreBoardV2/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const(
	push string = "$push"
	pull string = "$pull"
)

type LocationRepository interface{
	AddLocation()                    models.Result[models.Location]
	GetLocations()                   models.Result[models.Location]
	GetLocation(locationName string) models.Result[models.Location]
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
