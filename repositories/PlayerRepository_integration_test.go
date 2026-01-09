package repositories
// import (
// 	"JeopardyScoreBoardV2/models"
// 	"net/http"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestAddValidPlayer_Ok(t *testing.T) {
// 	playerRepository := GetSqlPlayerRepository(db)
// 	player := models.Player{PlayerName: "Darien Miller"}

// 	result := playerRepository.AddPlayerToLocation("Elmwood", player)

// 	assert.Equal(t, nil, result.Err)
// 	assert.Equal(t, http.StatusOK, result.StatusCode)
// 	assert.Equal(t, player.PlayerName, result.ResultData.PlayerName)

// 	//Verify that the player was inserted into the database
// 	allPlayers := playerRepository.GetAllPlayersFromAllLocations()
// 	playerInserted := allPlayers.ResultData[0]

// 	assert.Equal(t, player.PlayerName, playerInserted.PlayerName)
// }