package repositories

import(
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LocationRepository interface{

}

type MongoLocationRepository struct{
	locationCollection *mongo.Collection
}

func GetMongoLocationCollection(newMongoCollection *mongo.Collection) *MongoLocationRepository{
	return &MongoLocationRepository{ locationCollection: newMongoCollection }
}