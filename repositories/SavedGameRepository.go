package repositories

import (
	"JeopardyScoreBoardV2/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SavedGameRepository interface {
	GetAllSavedGames(req *http.Request) models.Result[[]models.SavedGame]
	GetAllSavedGamesFromLocation(req *http.Request, locationName string) models.Result[[]models.SavedGame]
	DeleteSavedGame(req *http.Request, savedGameId primitive.ObjectID) models.Result[*mongo.DeleteResult]
	AddSavedGame(req *http.Request, savedGame models.SavedGame) models.Result[models.SavedGame]
}