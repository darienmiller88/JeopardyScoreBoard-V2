package repositories

import (
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

//Receive a new instance of Location repository using postgres as the database. 
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
func (s *sqlLocationRepository) GetLocation(locationName string) models.Result[string]{
	location := ""
	
	if err := s.db.Get(&location, constants.GetLocation, location); err != nil{
		return getResult(err, http.StatusInternalServerError, "")
	}

	return getResult(nil, http.StatusOK, location)
}

//Helper function to allow repos to send result payloads with less text.
func getResult[T any](err error, statusCode int, payload T) models.Result[T] {
	return models.Result[T]{
		StatusCode: statusCode,
		Err: err,
		ResultData: payload,
	}
}