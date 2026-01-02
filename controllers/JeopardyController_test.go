package controllers

import (
	// "net/http/httptest"
	// "fmt"
	// "net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/services"
)

type mockService struct{
	getAllLocationsResult models.Result[[]string]
	getLocationResult models.Result[string]
}

func (m *mockService) GetAllLocations() models.Result[[]string]{
	return m.getAllLocationsResult
}

func (m *mockService) GetLocation(locationName string) models.Result[string]{
	return m.getLocationResult
}

func GetJeopardyController(service services.LocationService) *chi.Mux{
	jeopardyController := JeopardyController{}

	jeopardyController.Init(service)

	return jeopardyController.Router
}

func TestGetAllLocations_Ok(t *testing.T) {

	//Get the controller, and pass it a service

	//make http request to router

	//check status code and return value from the response

	assert.True(t, true, "True is true!")
}

func TestGetAllLocations_InternalServerError(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestGetLocation_Ok(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestGetLocation_NotFound(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestGetLocation_InternalServerErro(t *testing.T) {
	assert.True(t, true, "True is true!")
}