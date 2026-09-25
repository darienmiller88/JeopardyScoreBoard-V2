package repositories

import (
	"database/sql"
	"fmt"
	"net/http"

	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/encryption"
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/utils"

	"github.com/jmoiron/sqlx"
)

//LocationRepository interface to allow mocking when testing the service. The test can provide the service
//a dummy implementation
type LocationRepository interface{

	// GetLocation retrieves a single location name
	GetLocation(locationName string)  models.Result[string]

	// GetLocationById retrieves a single location by its ID
	GetLocationById(locationId int)   models.Result[models.Location]

	// GetAllLocations retrieves all location names in a string slice
	GetAllLocations()                 models.Result[[]string]
}

type locationRepository struct{
	db *sqlx.DB
	encryptionService *encryption.EncryptionService
}

//Receive a new instance of Location repository using postgres as the database. 
func NewLocationRepository(newDb *sqlx.DB, encryptionService *encryption.EncryptionService) LocationRepository{
	return &locationRepository{ 
		db: newDb, 
		encryptionService: encryptionService,
	}
}

//Retrieve all Locations from database
func (s *locationRepository) GetAllLocations() models.Result[[]string]{
	locations := []string{}
	
	if err := s.db.Select(&locations, constants.GetAllLocations); err != nil{
		return utils.GetResult(err, http.StatusInternalServerError, []string{})
	}

	return utils.GetResult(nil, http.StatusOK, locations)
}

//Get one location from the database
func (s *locationRepository) GetLocation(locationName string) models.Result[string]{
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
func (s *locationRepository) GetLocationById(locationId int) models.Result[models.Location]{
	location := models.Location{}
	
	if err := s.db.Get(&location, constants.GetLocationById, locationId); err != nil{
		if err == sql.ErrNoRows {
			return utils.GetResult(fmt.Errorf("No location found with id %d", locationId), http.StatusNotFound, location)	
		}

		return utils.GetResult(err, http.StatusInternalServerError, location)
	}

	return utils.GetResult(nil, http.StatusOK, location)
}