package repositories

import (
	"JeopardyScoreBoardV2/models"
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const(
	push string = "$push"
	pull string = "$pull"
)

type PlayerCardRepository interface {
	UpdatePlayerName(ctx context.Context, locationName string, oldPlayerName string, newPlayerName string) models.Result[*mongo.UpdateResult]
	AddPlayerToLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]
	RemovePlayerFromLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]
}

type MongoPlayerCardRepository struct{
	locationCollection *mongo.Collection
}

//Function to update a players name for a given location.
func (m *MongoPlayerCardRepository) UpdatePlayerName(ctx context.Context, locationName string, oldPlayerName string, newPlayerName string) models.Result[*mongo.UpdateResult]{
	//Filter first for the location, and then for the player in the "users" array field.
	filter := bson.M{ "location_name": locationName, "users.name": oldPlayerName }
	update := bson.M{ "$set": bson.M{ "users.$.name": newPlayerName } }

	//Using the above filters, find the specific user in the array at the specific location, and change their name.
	updateOneResult, err := m.locationCollection.UpdateOne(ctx, filter, update)

	if err != nil{
		return models.Result[*mongo.UpdateResult]{ Err: err, StatusCode: http.StatusInternalServerError }
	}

	return models.Result[*mongo.UpdateResult]{ ResultData: updateOneResult, StatusCode: http.StatusOK }
}

//Add a single player to a given location.
func (m *MongoPlayerCardRepository) AddPlayerToLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]{
	return m.updateHelper(ctx, push, locationName, playerName)
}

//Remove a single player from a given location.
func (m *MongoPlayerCardRepository) RemovePlayerFromLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]{
	return m.updateHelper(ctx, pull, locationName, playerName)
}

func (m *MongoPlayerCardRepository) updateHelper(ctx context.Context, mongoOperator string, locationName string, playerName string) models.Result[*mongo.UpdateResult]{
	filter := bson.M{ "location_name": locationName }
	update := bson.M{ mongoOperator: bson.M{"users": bson.M{"name": playerName} } }
 	updateOneResult, err := m.locationCollection.UpdateOne(ctx, filter, update)

	if err != nil{		
		return models.Result[*mongo.UpdateResult]{Err: err, StatusCode: http.StatusInternalServerError}
	}

	//If there were no documents modified, return an error and 404 signaling this.
	if updateOneResult.ModifiedCount == 0 {
		return models.Result[*mongo.UpdateResult]{
			Err: fmt.Errorf("no location \"%s\" found", locationName),
			StatusCode: http.StatusNotFound,
		}
	}

	return models.Result[*mongo.UpdateResult]{ ResultData: updateOneResult, StatusCode: http.StatusOK }
}