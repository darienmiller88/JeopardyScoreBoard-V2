 package services

// import (
// 	"net/http"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"

// 	"JeopardyScoreBoardV2/database"
// 	"JeopardyScoreBoardV2/models"
// )

// //This is a comment from my temp laptop to test commit success.
// func GetAllSavedGames(req *http.Request) models.Result[[]models.SavedGame] {
// 	findResult, err := database.GetSavedGamesCollections().Find(req.Context(), bson.D{})
	
// 	//Initialize the result after trying to find all current saved games.
// 	savedGamesResult := models.Result[[]models.SavedGame]{
// 		StatusCode: http.StatusInternalServerError,
// 		Err: err,
// 	}

// 	if err != nil {
// 		return savedGamesResult
// 	}

// 	savedGames := []models.SavedGame{}

// 	//Unmarshall the mongo cursor into the array of saved games.
// 	if err := findResult.All(req.Context(), &savedGames); err != nil{
// 		savedGamesResult.Err = err
// 		savedGamesResult.StatusCode = http.StatusInternalServerError
		
// 		return savedGamesResult
// 	}

// 	savedGamesResult.ResultData = savedGames
// 	savedGamesResult.StatusCode = http.StatusOK

// 	return savedGamesResult
// }