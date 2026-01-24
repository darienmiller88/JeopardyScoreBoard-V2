package services

import (
	"JeopardyScoreBoardV2/models"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLocationRepository struct{
	locationResult      models.Result[string]
	locationsResult     models.Result[[]string]
}

func (m *mockLocationRepository) GetAllLocations() models.Result[[]string]{
	return m.locationsResult
}

func (m *mockLocationRepository) GetLocation(locationName string) models.Result[string]{
	return m.locationResult
}

func (m *mockLocationRepository) IsLocationIdValid(locationId int) models.Result[string]{
	return m.locationResult
}

func TestGetAllLocations_Service_Ok(t *testing.T) {
	mockRepo := &mockLocationRepository{
		locationsResult: models.Result[[]string]{
			StatusCode: http.StatusOK,
			ResultData: []string{"Elmwood", "Flushing", "Pelham Bay"},
		},
	}

	service := &LocationServiceImpl{ Repository: mockRepo }
	result := service.GetAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 3)
	assert.Equal(t, "Elmwood", result.ResultData[0])
}

func TestGetAllLocations_Service_Error(t *testing.T) {
	mockRepo := &mockLocationRepository{
		locationsResult: models.Result[[]string]{
			Err:        fmt.Errorf("db down"),
			StatusCode: http.StatusInternalServerError,
		},
	}

	service := &LocationServiceImpl{ Repository: mockRepo }
	result := service.GetAllLocations()

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
}


func TestGetLocation_Service_Ok(t *testing.T) {
	location := "Elmwood"
	mockRepo := &mockLocationRepository{
		locationResult: models.Result[string]{
			StatusCode: http.StatusOK,
			ResultData: location,
		},
	}

	service := &LocationServiceImpl{Repository: mockRepo}
	result := service.GetLocation(location)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, location, result.ResultData)
}

func TestGetLocation_Service_Error(t *testing.T) {
	mockRepo := &mockLocationRepository{
		locationResult: models.Result[string]{
			Err:        fmt.Errorf("location not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &LocationServiceImpl{Repository: mockRepo}
	result := service.GetLocation("FakeTown")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Empty(t, result.ResultData)
}
