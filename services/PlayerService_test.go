package services

import (
	"net/http"
	"testing"

	"JeopardyScoreBoardV2/models"
)

type mockPlayerRepository struct{
	playerResult  models.Result[models.Player]
	playersResult models.Result[[]models.Player]
}

func (m *mockPlayerRepository) UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]{
	return m.playerResult
}
	
func (m *mockPlayerRepository) AddPlayerToLocation(locationName string, playerName string)  models.Result[models.Player]{
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
	mockRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{PlayerName: "Jane Doe"},
			StatusCode: http.StatusCreated,
		},
	}

	service := &PlayerServiceImpl{Repository: mockRepo}

	result := service.AddPlayerToLocation("Elmwood", "Jane Doe")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, "Jane Doe", result.Data.PlayerName)
}

func TestAddPlayer_PlayerNameTaken(t *testing.T){
	
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