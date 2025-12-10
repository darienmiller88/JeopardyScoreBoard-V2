package repositories

import (
	"JeopardyScoreBoardV2/models"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PlayerCardRepository interface {
	UpdatePlayerName(ctx context.Context, locationName string, oldPlayerName string, newPlayerName string) models.Result[*mongo.UpdateResult]
	AddPlayerToLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]
	RemovePlayerFromLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]
}

type MongoPlayerCardRepository struct{
	locationCollection *mongo.Collection
}

func (m *MongoPlayerCardRepository) UpdatePlayerName(ctx context.Context, locationName string, oldPlayerName string, newPlayerName string) models.Result[*mongo.UpdateResult]{
	updateUserResult := models.Result[*mongo.UpdateResult]{
		StatusCode: http.StatusInternalServerError,
	}
	
	//Filter first for the location, and then for the player in the "users" array field. Target the "name" 
	//field for each player field in the location collection
	filter := bson.M{ "location_name": locationName, "users.name": oldPlayerName }
	update := bson.M{ "$set": bson.M{ "users.$.name": newPlayerName } }

	//Using the above filters, find the specific user in the array at the specific location, and change their name.
	updateOneResult, err := m.locationCollection.UpdateOne(ctx, filter, update)

	if err != nil{
		updateUserResult.Err = err

		return updateUserResult
	}

	//If the update was succesful, return an Ok(200) and the update return data.
	updateUserResult.ResultData = updateOneResult
	updateUserResult.StatusCode = http.StatusOK

	return updateUserResult
}

func (m *MongoPlayerCardRepository) AddPlayerToLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]{
	return models.Result[*mongo.UpdateResult]{}
}

func (m *MongoPlayerCardRepository) RemovePlayerFromLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]{
	return models.Result[*mongo.UpdateResult]{}
}