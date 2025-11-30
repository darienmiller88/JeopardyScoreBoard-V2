package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client
const (
	databaseName string = "AdaptDB"
	locationsCollection string = "locations"
	savedGamesCollection string = "saved_games"
)

func Init() {
	var err error

	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}	

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func GetDB() *mongo.Client{
	return client
}

//Retrieve the "locations" collection from the database.
func GetLocationsCollection() *mongo.Collection {
	return client.Database(databaseName).Collection(locationsCollection)
}

//Retrieve the "saved_games" collection from the database.
func GetSavedGamesCollections() *mongo.Collection {
	return client.Database(databaseName).Collection(savedGamesCollection)
}

func DisconnectClient(){
	if err := client.Disconnect(context.TODO()); err != nil{
		panic(err)
	}
}