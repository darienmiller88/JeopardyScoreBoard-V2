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
	playerRepository := GetSqlPlayerRepository(db, nil)
	player := models.Player{PlayerName: "Darien Miller"}

	result := playerRepository.AddPlayerToLocation("Elmwood", player)

	assert.Equal(t, nil, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, player.PlayerName, result.ResultData.PlayerName)

	//Verify that the player was inserted into the database
	allPlayers := playerRepository.GetAllPlayersFromAllLocations()
	playerInserted := allPlayers.ResultData[len(allPlayers.ResultData) - 1]

	assert.Equal(t, player.PlayerName, playerInserted.PlayerName)
}

// func TestAddPlayerToLocation_IntegrationTest_InvalidLocation(t *testing.T) {
//     playerRepository := GetSqlPlayerRepository(db, nil)
// 	player := models.Player{PlayerName: "Darien Miller"}

// 	result := playerRepository.AddPlayerToLocation("FakeLocation", player)

// 	require.Error(t, result.Err)
// 	assert.Equal(t, http.StatusNotFound, result.StatusCode)
// }






/////////////////////////
//READ/GET tests
////////////////////////

func TestGetPlayersFromLocation_IntegrationTest_Ok(t *testing.T) {
 	playerRepository := GetSqlPlayerRepository(db, nil)
	result := playerRepository.GetPlayersFromLocation("Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestGetPlayersFromLocation_IntegrationTest_InvalidLocation(t *testing.T) {
 	playerRepository := GetSqlPlayerRepository(db, nil)
	result := playerRepository.GetPlayersFromLocation("FakeLocation")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, 0, len(result.ResultData))
}

func TestGetAllPlayersFromAllLocations_IntegrationTest_Ok(t *testing.T) {
 	playerRepository := GetSqlPlayerRepository(db, nil)
	result := playerRepository.GetAllPlayersFromAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}





/////////////////////////
//UPDATE/PUT tests
////////////////////////

func TestUpdatePlayerName_IntegrationTest_Happy(t *testing.T){
 	playerRepository := GetSqlPlayerRepository(db, nil)
	playerName := "player name"
	result := playerRepository.AddPlayerToLocation("Elmwood", models.Player{ PlayerName: playerName })

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)

	newName := "new name"
	updateResult := playerRepository.UpdatePlayerName(playerName, newName, "Elmwood")

	//Check ti see if the player's name was updated correctly
	require.NoError(t, updateResult.Err)
	assert.Equal(t, http.StatusOK, updateResult.StatusCode)
	assert.Equal(t, newName, updateResult.ResultData.PlayerName)
}

func TestUpdatePlayerName_IntegrationTest_PlayerNotFound_Happy(t *testing.T){
 	playerRepository := GetSqlPlayerRepository(db, nil)
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
    repo := GetSqlPlayerRepository(db, nil)

    // Insert a player
    _, err := db.Exec(`
        INSERT INTO players (player_name, location_id)
        VALUES ('DeleteMe', (SELECT id FROM locations WHERE location_name='Elmwood'))
    `)
    require.NoError(t, err)

    // Delete it
    result := repo.RemovePlayer("DeleteMe", "Elmwood")

    require.NoError(t, result.Err)
    assert.Equal(t, http.StatusOK, result.StatusCode)
    assert.Equal(t, "DeleteMe", result.ResultData.PlayerName)

    // Check to see if the player was removed
    var count int

    err = db.Get(&count, `SELECT COUNT(*) FROM players WHERE player_name='DeleteMe'`)

    require.NoError(t, err)
    assert.Equal(t, 0, count)
}

func TestRemovePlayer_IntegrationTest_PlayerNotFound(t *testing.T) {
    repo := GetSqlPlayerRepository(db, nil)

    // Make sure the table is empty
    // _, err := db.Exec(`DELETE FROM players`)
    // require.NoError(t, err)

    result := repo.RemovePlayer("NoName", "blag")

    require.Error(t, result.Err)
    assert.Equal(t, http.StatusNotFound, result.StatusCode)
    assert.Contains(t, result.Err.Error(), "could not find player")
}
