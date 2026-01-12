package services

import (
	"fmt"
	"net/http"
	"testing"

	"JeopardyScoreBoardV2/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlayerRepository struct{
	playerResult  models.Result[models.Player]
	playersResult models.Result[[]models.Player]
}

func (m *mockPlayerRepository) UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]{
	return m.playerResult
}
	
func (m *mockPlayerRepository) AddPlayerToLocation(locationName string, player models.Player)  models.Result[models.Player]{
	return m.playerResult
}

func (m *mockPlayerRepository) RemovePlayer(playerName string) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerRepository) GetPlayersFromLocation(locationName string) models.Result[[]models.Player] {
	return m.playersResult
}

func (m *mockPlayerRepository) GetAllPlayersFromAllLocations() models.Result[[]models.Player] {
	return m.playersResult
}


////////////////////
// CREATE/POST tests
////////////////////

func TestAddPlayer_Ok(t *testing.T){
	validPlayerName := "Jane Doe" //Valid name, two parts and vlaid length
	mockRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{ PlayerName: validPlayerName },
			StatusCode: http.StatusCreated,
		},
	}

	service := &PlayerServiceImpl{ Repository: mockRepo }
	result := service.AddPlayerToLocation("Elmwood", validPlayerName)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, "Jane Doe", result.ResultData.PlayerName)
}

func TestAddPlayer_NameTooShort(t *testing.T) {
	mockRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{ Repository: mockRepo }

	result := service.AddPlayerToLocation("Elmwood", "Joe")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

func TestAddPlayer_NameTooLong(t *testing.T) {
	mockRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{ Repository: mockRepo }

	result := service.AddPlayerToLocation("Elmwood", "Joedcsxevrgvfsxergtdwxertgfwsxdgtrvwsxdertgvcevrtbgfcdw")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

func TestAddPlayer_NameMustHaveTwoParts(t *testing.T) {
	mockRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{ Repository: mockRepo }

	result := service.AddPlayerToLocation("Elmwood", "Cheryl")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "two parts")

	result = service.AddPlayerToLocation("Elmwood", "   This name has more than two parts   ")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "two parts")
}




////////////////////
// READ/GET tests
////////////////////


func TestGetAllPlayersFromAllLocations_Ok(t *testing.T) {
	mockRepo := &mockPlayerRepository{
		playersResult: models.Result[[]models.Player]{
			ResultData: []models.Player{
				{ PlayerName: "Jane Doe" },
				{ PlayerName: "John Smith" },
			},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{ Repository: mockRepo }
	result := service.GetAllPlayersFromAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 2)
}


func TestGetPlayersFromLocation_Ok(t *testing.T) {
	mockRepo := &mockPlayerRepository{
		playersResult: models.Result[[]models.Player]{
			ResultData: []models.Player{
				{ PlayerName: "Jane Doe" },
			},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{ Repository: mockRepo }
	result := service.GetPlayersFromLocation("Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 1)
}

func TestGetPlayersFromLocation_RepoError(t *testing.T) {
	mockRepo := &mockPlayerRepository{
		playersResult: models.Result[[]models.Player]{
			Err:        fmt.Errorf("db exploded"),
			StatusCode: http.StatusInternalServerError,
		},
	}

	service := &PlayerServiceImpl{ Repository: mockRepo }
	result := service.GetPlayersFromLocation("Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}


////////////////////
// UPDATE/PUT tests
////////////////////

func TestUpdatePlayerName_Service_Ok(t *testing.T) {
	mockRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{ PlayerName: "Bob Melendez" },
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{ Repository: mockRepo }
	result := service.UpdatePlayerName("Alice Twilight", "Bob Melendez")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "Bob Melendez", result.ResultData.PlayerName)
}

func TestUpdatePlayerName_Service_InvalidNewName(t *testing.T) {
	mockRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{ Repository: mockRepo }

	result := service.UpdatePlayerName("Alice Twilight", "Bob") // too short (3 chars)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

func TestUpdatePlayerName_Service_NameMustHaveTwoParts(t *testing.T) {
	mockRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{Repository: mockRepo}

	result := service.UpdatePlayerName("Alice Twilight", "Margaret")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

func TestUpdatePlayerName_Service_SameName(t *testing.T) {
	mockRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{ Repository: mockRepo }

	result := service.UpdatePlayerName("Jane Doe", "Jane Doe")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "must be different")
}

func TestUpdatePlayerName_Service_RepoError(t *testing.T) {
	mockRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("db down"),
			StatusCode: http.StatusInternalServerError,
		},
	}

	service := &PlayerServiceImpl{ Repository: mockRepo }
	result := service.UpdatePlayerName("Alice", "Bob Melendez")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}


////////////////////
// DESTROY/DELETE tests
////////////////////