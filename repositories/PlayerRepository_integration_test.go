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
	playerInserted := allPlayers.ResultData[len(allPlayers.ResultData) - 1]

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

func TestGetPlayersFromLocation_IntegrationTest_Ok(t *testing.T) {
 	playerRepository := GetSqlPlayerRepository(db)
	result := playerRepository.GetPlayersFromLocation("Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestGetPlayersFromLocation_IntegrationTest_InvalidLocation(t *testing.T) {
 	playerRepository := GetSqlPlayerRepository(db)
	result := playerRepository.GetPlayersFromLocation("FakeLocation")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, 0, len(result.ResultData))
}

func TestGetAllPlayersFromAllLocations_IntegrationTest_Ok(t *testing.T) {
 	playerRepository := GetSqlPlayerRepository(db)
	result := playerRepository.GetAllPlayersFromAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

/////////////////////////
//UPDATE/PUT tests
////////////////////////

func TestUpdatePlayerName_IntegrationTest_Ok(t *testing.T){
 	playerRepository := GetSqlPlayerRepository(db)
	playerName := "player name"
	result := playerRepository.AddPlayerToLocation("Elmwood", models.Player{ PlayerName: playerName })

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)

	newName := "new name"
	updateResult := playerRepository.UpdatePlayerName(playerName, newName)

	//Check ti see if the player's name was updated correctly
	require.NoError(t, updateResult.Err)
	assert.Equal(t, http.StatusOK, updateResult.StatusCode)
	assert.Equal(t, newName, updateResult.ResultData.PlayerName)
}

func TestUpdatePlayerName_IntegrationTest_NewNameTaken(t *testing.T) {
	repo := GetSqlPlayerRepository(db)

    // Create two players in same location
    _, err := db.Exec(`
        INSERT INTO players (player_name, location_id)
        VALUES
          ('Alice', (SELECT id FROM locations WHERE location_name='Elmwood')),
          ('Bob',   (SELECT id FROM locations WHERE location_name='Elmwood'))
    `)
    require.NoError(t, err)

    // Try to rename Alice -> Bob (duplicate name)
    result := repo.UpdatePlayerName("Alice", "Bob")

    require.Error(t, result.Err)
    assert.Equal(t, http.StatusConflict, result.StatusCode)
    assert.Contains(t, result.Err.Error(), "already exists")
}

func TestUpdatePlayerName_IntegrationTest_OldNameNotFound(t *testing.T) {
	repo := GetSqlPlayerRepository(db)

    // Ensure table is empty 
    _, err := db.Exec(`DELETE FROM players`)
    require.NoError(t, err)

    result := repo.UpdatePlayerName("Ghost", "NewName")

    require.Error(t, result.Err)
    assert.Equal(t, http.StatusNotFound, result.StatusCode)
    assert.Contains(t, result.Err.Error(), "could not find player")
}


/////////////////////////
//DESTROY/DELETE tests
////////////////////////

func TestRemovePlayer_IntegrationTest_Ok(t *testing.T) {
    repo := GetSqlPlayerRepository(db)

    // Insert a player
    _, err := db.Exec(`
        INSERT INTO players (player_name, location_id)
        VALUES ('DeleteMe', (SELECT id FROM locations WHERE location_name='Elmwood'))
    `)
    require.NoError(t, err)

    // Delete it
    result := repo.RemovePlayer("DeleteMe")

    require.NoError(t, result.Err)
    assert.Equal(t, http.StatusOK, result.StatusCode)
    assert.Equal(t, "DeleteMe", result.ResultData.PlayerName)

    // Check to see if the player was removed
    var count int

    err = db.Get(&count, `SELECT COUNT(*) FROM players WHERE player_name='DeleteMe'`)
	
    require.NoError(t, err)
    assert.Equal(t, 0, count)
}

