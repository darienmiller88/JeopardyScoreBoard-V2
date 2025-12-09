package repositories

import(
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/database"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
func GetMongoLocationCollection() *MongoLocationRepository{
	return &MongoLocationRepository{ locationCollection: database.GetLocationsCollection().Clone() }
}

func (m *MongoLocationRepository) AddLocation() models.Result[models.Location]{
	return models.Result[models.Location]{}
}

func (m *MongoLocationRepository) GetLocations() models.Result[models.Location]{
	return models.Result[models.Location]{}
}