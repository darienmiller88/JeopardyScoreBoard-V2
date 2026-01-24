package repositories

import (
	"database/sql"
	"fmt"
	"net/http"

	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/utils"

	"github.com/jmoiron/sqlx"
)

//LocationRepository interface to allow mocking when testing the service. The test can provide the service
//a dummy implementation
type LocationRepository interface{
	GetLocation(locationName string)  models.Result[string]
	GetAllLocations()                 models.Result[[]string]
	IsLocationIdValid(locationId int) models.Result[string]
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
		return utils.GetResult(err, http.StatusInternalServerError, []string{})
	}

	return utils.GetResult(nil, http.StatusOK, locations)
}

//Get one location from the database
func (s *sqlLocationRepository) GetLocation(locationName string) models.Result[string]{
	location := ""
	
	if err := s.db.Get(&location, constants.GetLocation, locationName); err != nil{
		if err == sql.ErrNoRows {
			return utils.GetResult(fmt.Errorf("No location found with name %s", locationName), http.StatusNotFound, "")	
		}

		return utils.GetResult(err, http.StatusInternalServerError, "")
	}

	return utils.GetResult(nil, http.StatusOK, location)
}

//Get one location from the database
func (s *sqlLocationRepository) IsLocationIdValid(locationId int) models.Result[string]{
	location := ""
	
	if err := s.db.Get(&location, constants.GetLocationById, locationId); err != nil{
		if err == sql.ErrNoRows {
			return utils.GetResult(fmt.Errorf("No location with id %d", locationId), http.StatusNotFound, "")	
		}

		return utils.GetResult(err, http.StatusInternalServerError, "")
	}

	return utils.GetResult(nil, http.StatusOK, location)
}
