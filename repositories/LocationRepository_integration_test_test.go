package repositories

import (
	"net/http"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test for getting all 8 locations
//EXPECTED PAYLOAD:       All 8 locations
//EXPECTED STATUS CODE:   200
//EXPECTED ERROR MESSAGE: nil          
func TestGetAllLocations_IntegrationTest_Ok(t *testing.T){
	locationRepository := GetSqlLocationRepository(db)	

 	result := locationRepository.GetAllLocations()
	
	//As part of the migrations, 8 locations are inserted into the database. We need to ensure all 8 are returned.
	expectedNumLocations := 8

	assert.Equal(t, nil, result.Err)
	assert.Equal(t, expectedNumLocations, len(result.ResultData)) 
}

//Test for when a location to be searched does exist.
//EXPECTED PAYLOAD:       location_name
//EXPECTED STATUS CODE:   200
//EXPECTED ERROR MESSAGE: nil            
func TestGetLocation_IntegrationTest_Ok(t *testing.T){
	locationRepository := GetSqlLocationRepository(db)	

	//"Elmwood" is one of the locations 
 	result := locationRepository.GetLocation("Elmwood")
	
	assert.Equal(t, nil, result.Err)
	assert.Equal(t, "Elmwood", result.ResultData) 
}

//Test for when a location to be searched does not exist
//EXPECTED PAYLOAD:       ""
//EXPECTED STATUS CODE:   404
//EXPECTED ERROR MESSAGE: No location found with name location_name            
func TestGetLocation_IntegrationTest_NotFound(t *testing.T){
	locationRepository := GetSqlLocationRepository(db)	

	locationToCheck := "FakeLocation"
 	result := locationRepository.GetLocation(locationToCheck)
	
	assert.Equal(t, fmt.Errorf("No location found with name %s", locationToCheck), result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Equal(t, "", result.ResultData) 
}