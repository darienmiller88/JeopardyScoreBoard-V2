package repositories

import(
	"JeopardyScoreBoardV2/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LocationRepository interface{
	AddLocation() models.Result[models.Location]
	GetLocation() models.Result[models.Location]
}

type MongoLocationRepository struct{
	locationCollection *mongo.Collection
}

func GetMongoLocationCollection(newMongoCollection *mongo.Collection) *MongoLocationRepository{
	return &MongoLocationRepository{ locationCollection: newMongoCollection }
}