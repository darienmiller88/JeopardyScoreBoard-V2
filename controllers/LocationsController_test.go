package controllers

import (
	// "fmt"
	"errors"
	"io"
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

func getLocationsController(service services.LocationService) *chi.Mux{
	jeopardyController := LocationsController{}

	jeopardyController.Init(service)

	return jeopardyController.Router
}

func getResponseFromController(router *chi.Mux, route string, method string, body io.Reader) (*http.Response, *httptest.ResponseRecorder){
	req := httptest.NewRequest(method, route, body)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	
	return res.Result(), res
}

func TestGetAllLocations_Ok(t *testing.T) {
	//Create a mock service that returns a list of locations alongside a 200. This will simulate the Jeopardy
	//controller calling the location service, which retrieves the ADAPT locations from the database through
	//the repository
	router := getLocationsController(&mockService{ 
		getAllLocationsResult: models.Result[[]string]{
			Err: nil,
			StatusCode: http.StatusOK,
			ResultData: []string{"Pelham Bay", "Elmwood"},
		},
	})

	result, res := getResponseFromController(router, "/", http.MethodGet)

	//check status code and return value from the response, and see if the status code is 200
	assert.Equal(t, http.StatusOK, result.StatusCode)

	//Check to see if the return body contains the following locations.
	assert.Contains(t, res.Body.String(), "Pelham Bay")
	assert.Contains(t, res.Body.String(), "Elmwood")
}

func TestGetAllLocations_InternalServerError(t *testing.T) {
	//Create a mock service that returns a 500, and an empty list. This will simulate a internal server
	//error when retirieving the locations
	router := getLocationsController(&mockService{ 
		getAllLocationsResult: models.Result[[]string]{
			Err: errors.New("Database error"),
			StatusCode: http.StatusInternalServerError,
		},
	})

	result, res := getResponseFromController(router, "/", http.MethodGet)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Contains(t, res.Body.String(), "Database error")
}

func TestGetLocation_Ok(t *testing.T) {
	//Create a mock service that returns a 200, and a location. 
	router := getLocationsController(&mockService{ 
		getLocationResult: models.Result[string]{
			Err: nil,
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	})

	result, res := getResponseFromController(router, "/Elmwood", http.MethodGet)
	
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Contains(t, res.Body.String(), "Elmwood")
}

func TestGetLocation_NotFound(t *testing.T) {
	//Create a mock service that returns a 404. This will simulate the database not finding a location.
	router := getLocationsController(&mockService{ 
		getLocationResult: models.Result[string]{
			Err: errors.New("Location not found"),
			StatusCode: http.StatusNotFound,
		},
	})

	result, res := getResponseFromController(router, "/fakelocation", http.MethodGet)

	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Contains(t, res.Body.String(), "Location not found")
}

func TestGetLocation_InternalServerError(t *testing.T) {
	//Create a mock service that returns a 500. This will simulate a internal server error when retirieving 
	//one location
	router := getLocationsController(&mockService{ 
		getLocationResult: models.Result[string]{
			Err: errors.New("Internal Server Error"),
			StatusCode: http.StatusInternalServerError,
		},
	})

	result, res := getResponseFromController(router, "/Elmwood", http.MethodGet)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Contains(t, res.Body.String(), "Internal Server Error")
}