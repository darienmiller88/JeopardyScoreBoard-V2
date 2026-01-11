package repositories

import (
	"JeopardyScoreBoardV2/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/////////////////////////
//CREATE/POST tests
////////////////////////
func TestAddValidPlayer_IntegrationTest_Ok(t *testing.T) {
	playerRepository := GetSqlPlayerRepository(db)
	player := models.Player{PlayerName: "Darien Miller"}

	result := playerRepository.AddPlayerToLocation("Elmwood", player)

	assert.Equal(t, nil, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, player.PlayerName, result.ResultData.PlayerName)

	//Verify that the player was inserted into the database
	allPlayers := playerRepository.GetAllPlayersFromAllLocations()
	playerInserted := allPlayers.ResultData[0]

	assert.Equal(t, player.PlayerName, playerInserted.PlayerName)
}

func TestAddPlayerToLocation_IntegrationTest_InvalidLocation(t *testing.T) {
    playerRepository := GetSqlPlayerRepository(db)
	player := models.Player{PlayerName: "Darien Miller"}

	result := playerRepository.AddPlayerToLocation("FakeLocation", player)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

func TestAddPlayerToLocation_IntegrationTest_PlayerNameTaken(t *testing.T) {
    playerRepository := GetSqlPlayerRepository(db)
	player := models.Player{PlayerName: "Darien Miller"}

	result := playerRepository.AddPlayerToLocation("Elmwood", player)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
	assert.Equal(t, models.Player{}, result.ResultData)
}


/////////////////////////
//READ/GET tests
////////////////////////





/////////////////////////
//UPDATE/PUT tests
////////////////////////





/////////////////////////
//DESTROY/DELETE tests
////////////////////////