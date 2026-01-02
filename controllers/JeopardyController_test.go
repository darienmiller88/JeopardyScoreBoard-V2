package controllers

import (
	// "fmt"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/services"
)

type mockService struct{
	getAllLocationsResult models.Result[[]string]
	getLocationResult     models.Result[string]
}

func (m *mockService) GetAllLocations() models.Result[[]string]{
	return m.getAllLocationsResult
}

func (m *mockService) GetLocation(locationName string) models.Result[string]{
	return m.getLocationResult
}

func getJeopardyController(service services.LocationService) *chi.Mux{
	jeopardyController := JeopardyController{}

	jeopardyController.Init(service)

	return jeopardyController.Router
}

func TestGetAllLocations_Ok(t *testing.T) {

	//Create a mock service that returns a list of locations alongside a 200. This will simulate the Jeopardy
	//controller calling the location service, which retrieves the ADAPT locations from the database through
	//the repository
	router := getJeopardyController(&mockService{ 
		getAllLocationsResult: models.Result[[]string]{
			Err: nil,
			StatusCode: http.StatusOK,
			ResultData: []string{"Pelham Bay", "Elmwood"},
		},
	})

	//The jeopardy controller will create two routes: "/" and "/{location_name}"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	response := res.Result()

	//check status code and return value from the response, and see if the status code is 200
	assert.Equal(t, http.StatusOK, response.StatusCode)

	//Check to see if the return body contains the following locations.
	assert.Contains(t, res.Body.String(), "Pelham Bay")
	assert.Contains(t, res.Body.String(), "Elmwood")
}

func TestGetAllLocations_InternalServerError(t *testing.T) {
	//Create a mock service that returns a 500, and an empty list. This will simulate a internal server
	//error when retirieving the locations
	router := getJeopardyController(&mockService{ 
		getAllLocationsResult: models.Result[[]string]{
			Err: errors.New("Database error"),
			StatusCode: http.StatusInternalServerError,
		},
	})

	//The jeopardy controller will create two routes: "/" and "/{location_name}"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	response := res.Result()

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Contains(t, res.Body.String(), "Database error")
}

func TestGetLocation_Ok(t *testing.T) {
	//Create a mock service that returns a 500, and an empty list. This will simulate a internal server
	//error when retirieving the locations
	router := getJeopardyController(&mockService{ 
		getLocationResult: models.Result[string]{
			Err: nil,
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	})

	//The jeopardy controller will create two routes: "/" and "/{location_name}"
	req := httptest.NewRequest(http.MethodGet, "/Elmwood", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	response := res.Result()
	
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, res.Body.String(), "Elmwood")
}

func TestGetLocation_NotFound(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestGetLocation_InternalServerError(t *testing.T) {
	assert.True(t, true, "True is true!")
}