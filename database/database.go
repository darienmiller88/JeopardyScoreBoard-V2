package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// var client *mongo.Client
var db *sqlx.DB

func Init(){
	_db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))

	if err != nil{
		panic(err)
	}
	
	err = _db.Ping()
	
	if err != nil {
		fmt.Println("db connection fail:", err)
	}else{
		fmt.Println("Connection established! :)")
	}
	
	db = _db
}
	
func GetDB() *sqlx.DB{
	return db
}

func CloseSQLDB(){
	if err := db.Close(); err != nil{
		panic(err)
	}
}
	
// const (
// 	databaseName string = "AdaptDB"
// 	locationsCollection string = "locations"
// 	savedGamesCollection string = "saved_games"
// )
	
// func Init() {
// 	var err error

// 	client, err = mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := client.Ping(context.TODO(), nil); err != nil {
// 		panic(err)
// 	}	

// 	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
// }


//Retrieve the "locations" collection from the database.
// func GetLocationsCollection() *mongo.Collection {
// 	return client.Database(databaseName).Collection(locationsCollection)
// }

// //Retrieve the "saved_games" collection from the database.
// func GetSavedGamesCollections() *mongo.Collection {
// 	return client.Database(databaseName).Collection(savedGamesCollection)
// }

// func DisconnectClient(){
// 	if err := client.Disconnect(context.TODO()); err != nil{
// 		panic(err)
// 	}
// }
