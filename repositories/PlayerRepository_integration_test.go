package repositories

import (
	"JeopardyScoreBoardV2/encryption"
	"JeopardyScoreBoardV2/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ///////////////////////
// CREATE/POST tests
// //////////////////////
func TestAddValidPlayer_IntegrationTest_Ok(t *testing.T) {
	playerRepository := GetSqlPlayerRepository(db, es)
	player := models.Player{PlayerName: "Darien Miller"}

	result := playerRepository.AddPlayerToLocation("Elmwood", player)

	assert.Equal(t, nil, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)

	nameHash := encryption.NameHash(player.PlayerName)

	//Verify the player name hash 
	assert.Equal(t, nameHash, result.ResultData.PlayerNameHash)

	//Verify that the player was added to the database.
	hash := []byte{}
	db.Get(&hash, "SELECT player_name_hash FROM players WHERE id=$1", result.ResultData.ID)

	assert.Equal(t, nameHash, hash)
}


/////////////////////////
//READ/GET tests
////////////////////////

func TestGetPlayersFromLocation_IntegrationTest_Ok(t *testing.T) {
	playerRepository := GetSqlPlayerRepository(db, es)
	result := playerRepository.GetPlayersFromLocation("Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestGetPlayersFromLocation_IntegrationTest_InvalidLocation(t *testing.T) {
	playerRepository := GetSqlPlayerRepository(db, es)
	result := playerRepository.GetPlayersFromLocation("FakeLocation")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, 0, len(result.ResultData))
}

func TestGetAllPlayersFromAllLocations_IntegrationTest_Ok(t *testing.T) {
	playerRepository := GetSqlPlayerRepository(db, es)
	result := playerRepository.GetAllPlayersFromAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

/////////////////////////
//UPDATE/PUT tests
////////////////////////

func TestUpdatePlayerName_IntegrationTest_Happy(t *testing.T) {
	playerRepository := GetSqlPlayerRepository(db, es)
	playerName := "player name"
	result := playerRepository.AddPlayerToLocation("Elmwood", models.Player{PlayerName: playerName})

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)

	newName := "new name"
	updateResult := playerRepository.UpdatePlayerName(playerName, newName, "Elmwood")

	//Check ti see if the player's name was updated correctly
	require.NoError(t, updateResult.Err)
	assert.Equal(t, http.StatusOK, updateResult.StatusCode)
	assert.Equal(t, newName, updateResult.ResultData.PlayerName)
}

func TestUpdatePlayerName_IntegrationTest_PlayerNotFound_Happy(t *testing.T) {
	playerRepository := GetSqlPlayerRepository(db, es)
	updateResult := playerRepository.UpdatePlayerName("fakename", "newName", "Elmwood")

	//check to see if an error was returned, and a 404 was sent as well
	require.Error(t, updateResult.Err)
	assert.Equal(t, http.StatusNotFound, updateResult.StatusCode)
	assert.Contains(t, updateResult.Err.Error(), "could not find player")
}

/////////////////////////
//DESTROY/DELETE tests
////////////////////////

func TestRemovePlayer_IntegrationTest_Ok(t *testing.T) {
	repo := GetSqlPlayerRepository(db, es)
	player := models.Player{PlayerName: "fuuun nammmee"}

	// Insert a player
	result := repo.AddPlayerToLocation("Elmwood", player)
	require.NoError(t, result.Err)

	// Delete it
	result = repo.RemovePlayer(player.PlayerName, "Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)

	//Check to see if the player was found.
	result = repo.GetPlayerByName(player.PlayerName)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

func TestRemovePlayer_IntegrationTest_PlayerNotFound(t *testing.T) {
	repo := GetSqlPlayerRepository(db, es)

	// Make sure the table is empty
	// _, err := db.Exec(`DELETE FROM players`)
	// require.NoError(t, err)

	result := repo.RemovePlayer("NoName", "blag")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "could not find player")
}
