package repositories

import (
	"context"
	"net/http"
	
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/constants"

	"github.com/jmoiron/sqlx"
)

//LocationRepository interface to allow mocking when testing the service. The test can provide the service
//a dummy implementation
type LocationRepository interface{
	GetLocation(locationName string) models.Result[string]
	GetAllLocations()                models.Result[[]string]
}

type sqlLocationRepository struct{
	db *sqlx.DB
}

//Receive a new instance of Location repository using a mongo collection as the database. 
func GetSqlLocationRepository(newDb *sqlx.DB) *sqlLocationRepository{
	return &sqlLocationRepository{ db: newDb }
}

//Retrieve all Locations from database
func (s *sqlLocationRepository) GetAllLocations() models.Result[[]string]{
	locations := []string{}
	
	if err := s.db.Select(&locations, constants.GetAllLocations); err != nil{
		return getResult(err, http.StatusInternalServerError, []string{})
	}

	return getResult(nil, http.StatusOK, locations)
}

//Get one location from the database
func (s *sqlLocationRepository) GetLocation(ctx context.Context, locationName string) models.Result[string]{
	location := ""
	
	if err := s.db.Get(&location, constants.GetLocation, location); err != nil{
		return getResult(err, http.StatusInternalServerError, "")
	}

	return getResult(nil, http.StatusOK, location)
}

// //Return the array of players from a particular location
// func (m *MongoLocationRepository) GetPlayersFromLocation(ctx context.Context, locationName string) models.Result[[]models.PlayerCard]{
// 	locationResult := m.GetLocation(ctx, locationName)

// 	//Return the error from the above call if an error occurs
// 	if locationResult.Err != nil {
// 		return getResult(locationResult.Err, http.StatusInternalServerError, []models.PlayerCard{})
// 	}

// 	//IF not, return all of the players for for that location.
// 	return getResult(nil, http.StatusOK, locationResult.ResultData.Players)
// }

// //Return every single player at every single location.
// func (m *MongoLocationRepository) GetAllPlayersFromAllLocations(ctx context.Context) models.Result[[]models.PlayerCard]{
// 	locationsResult := m.GetAllLocations(ctx)	

// 	if locationsResult.Err != nil {
// 		return getResult(locationsResult.Err, locationsResult.StatusCode, []models.PlayerCard{})
// 	}

// 	players := []models.PlayerCard{}

// 	//Range over each location, and extract all of the players for each location.
// 	for _, location := range locationsResult.ResultData {
// 		players = append(players, location.Players...)
// 	}

// 	//Return the list of all players in the database.
// 	return getResult(nil, http.StatusOK, players)
// }

//Helper function to allow repos to send result payloads with less text.
func getResult[T any](err error, statusCode int, payload T) models.Result[T] {
	return models.Result[T]{
		StatusCode: statusCode,
		Err: err,
		ResultData: payload,
	}
}