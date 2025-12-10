package repositories

import (
	"JeopardyScoreBoardV2/models"
	"context"
	"fmt"
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

//Add a single player to a given location.
func (m *MongoPlayerCardRepository) AddPlayerToLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]{
	updateUserResult := models.Result[*mongo.UpdateResult]{
		StatusCode: http.StatusInternalServerError,
	}

	filter := bson.M{ "location_name": locationName }
	update := bson.M{ "$push": bson.M{"users": bson.M{"name": playerName} } }
 	updateOneResult, err := m.locationCollection.UpdateOne(ctx, filter, update)

	if err != nil{
		updateUserResult.Err = err
		
		return updateUserResult
	}

	if updateOneResult.ModifiedCount == 0 {
		updateUserResult.Err = fmt.Errorf("no location \"%s\" found", locationName)
		updateUserResult.StatusCode = http.StatusNotFound

		return updateUserResult
	}

	updateUserResult.ResultData = updateOneResult
	updateUserResult.StatusCode = http.StatusOK

	return updateUserResult
}

func (m *MongoPlayerCardRepository) RemovePlayerFromLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]{
	filter := bson.M{ "location_name": locationName }
	update := bson.M{ "$pull": bson.M{"users": bson.M{"name": playerName} } }
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