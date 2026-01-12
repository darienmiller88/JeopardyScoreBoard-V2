package services

import (
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
}




////////////////////
// READ/GET tests
////////////////////







////////////////////
// UPDATE/PUT tests
////////////////////







////////////////////
// DESTROY/DELETE tests
////////////////////