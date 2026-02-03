package repositories

import (
	"JeopardyScoreBoardV2/models"
	"database/sql"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

//=======================================
// GET / tests retrieving saved games
//======================================

// GetAllSavedGamesDB_Integration_Happy
// Verifies all seeded saved games are returned from the real test database.
func TestGetAllSavedGamesDB_Integration_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)
	result := repo.GetAllSavedGamesDB()

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)

	// From migrations: 8 total saved games
	require.GreaterOrEqual(t, len(result.ResultData), 8)

	// Spot check a known seeded value
	found := false
	for _, sg := range result.ResultData {
		if sg.TotalPoints == 1200 {
			found = true
			break
		}
	}
	require.True(t, found, "expected seeded saved game not found")
}

// GetAllSavedGamesDB_Integration_Unhappy_DBClosed
// Verifies error is returned when the DB connection is closed.
// func TestGetAllSavedGamesDB_Integration_Unhappy_DBClosed(t *testing.T) {
// 	badDB := db
// 	badDB.Close() // simulate catastrophic failure

// 	repo := GetSqlSavedGameRepository(badDB)

// 	result := repo.GetAllSavedGamesDB()

// 	require.NotNil(t, result.Err)
// 	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
// 	require.Empty(t, result.ResultData)
// }


// GetAllSavedGamesFromLocationDB_Integration_Happy_Elmwood
// Verifies Elmwood returns the 2 seeded saved games.
func TestGetAllSavedGamesFromLocationDB_Integration_Happy_Elmwood(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	result := repo.GetAllSavedGamesFromLocationDB("Elmwood")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)
}

// GetAllSavedGamesFromLocationDB_Integration_Happy_Lawrence
// Verifies Lawrence returns the 4 seeded saved games.
func TestGetAllSavedGamesFromLocationDB_Integration_Happy_Lawrence(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)
	result := repo.GetAllSavedGamesFromLocationDB("Lawrence")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 4)
}

// GetAllSavedGamesFromLocationDB_Integration_Happy_Flushing
// Verifies Flushing returns the 2 seeded saved games.
func TestGetAllSavedGamesFromLocationDB_Integration_Happy_Flushing(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	result := repo.GetAllSavedGamesFromLocationDB("Flushing")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)
}

// GetAllSavedGamesFromLocationDB_Integration_Unhappy_NoGames
// Verifies empty slice when location exists but has no saved games.
func TestGetAllSavedGamesFromLocationDB_Integration_Unhappy_NoGames(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)
	result := repo.GetAllSavedGamesFromLocationDB("Pelham Bay")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Empty(t, result.ResultData)
}

// GetAllSavedGamesFromLocationDB_Integration_Unhappy_BadLocation
// Verifies empty slice when location name does not exist.
func TestGetAllSavedGamesFromLocationDB_Integration_Unhappy_BadLocation(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	result := repo.GetAllSavedGamesFromLocationDB("NotARealLocation")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Empty(t, result.ResultData)
}

// GetAllSavedGamesFromLocationDB_Integration_Unhappy_DBClosed
// Verifies error when DB connection is closed.
// func TestGetAllSavedGamesFromLocationDB_Integration_Unhappy_DBClosed(t *testing.T) {
// 	badDB := db
// 	badDB.Close()

// 	repo := GetSqlSavedGameRepository(badDB)
// 	result := repo.GetAllSavedGamesFromLocationDB("Elmwood")

// 	require.NotNil(t, result.Err)
// 	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
// 	require.Empty(t, result.ResultData)
// }




//======================================
// DELETE / tests deleting saved games
//======================================

func getAnySavedGameID(t *testing.T) string {
	var id string
	err := db.Get(&id, "SELECT id FROM savedgames LIMIT 1")
	require.NoError(t, err)
	return id
}

func TestDeleteSavedGameDB_Integration_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	// Get a real saved game id
	id := getAnySavedGameID(t)

	// Count before
	var before int
	err := db.Get(&before, "SELECT COUNT(*) FROM savedgames")
	require.NoError(t, err)

	result := repo.DeleteSavedGameDB(id)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Contains(t, result.ResultData, id)

	// Count after
	var after int
	err = db.Get(&after, "SELECT COUNT(*) FROM savedgames")
	require.NoError(t, err)

	require.Equal(t, before-1, after)

	// Verify it is truly gone
	var exists int
	err = db.Get(&exists, "SELECT COUNT(*) FROM savedgames WHERE id=$1", id)
	require.NoError(t, err)
	require.Equal(t, 0, exists)
}

func TestDeleteSavedGameDB_Integration_Unhappy_NotFound(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)
	fakeID := "999"
	result := repo.DeleteSavedGameDB(fakeID)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusNotFound, result.StatusCode)
}

// func TestDeleteSavedGameDB_Integration_Unhappy_DBClosed(t *testing.T) {
// 	badDB := db
// 	badDB.Close()

// 	repo := GetSqlSavedGameRepository(badDB)
// 	result := repo.DeleteSavedGameDB("9999")

// 	require.NotNil(t, result.Err)
// 	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
// }




//=======================================
// INSERT - tests inserting saved games
//=======================================

func getLocationID(t *testing.T, name string) int {
	var id int
	err := db.Get(&id, "SELECT id FROM locations WHERE location_name=$1", name)
	require.NoError(t, err)
	return id
}

func getPlayerID(t *testing.T, name string) int {
	var id int
	err := db.Get(&id, "SELECT id FROM players WHERE player_name=$1", name)
	require.NoError(t, err)
	return id
}

func TestAddStandardSavedGame_Integration_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	locationID := getLocationID(t, "Elmwood")
	playerOneID := getPlayerID(t, "playerone")
	playerTwoID := getPlayerID(t, "playertwo")

	savedGame := models.SavedGame{
		TotalPoints:       900,
		AveragePoints:     300,
		WinningPlayerName: sql.NullString{String: "playerone"},
		WinningPlayerId:   sql.NullInt32{Int32: int32(playerOneID), Valid: true },
		LocationId:        locationID,
		Players: []models.Player{
			{ID: playerOneID, PlayerName: "playerone", Score: 500},
			{ID: playerTwoID, PlayerName: "playertwo", Score: 400},
		},
	}

	result := repo.addStandardSavedGame(savedGame)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.NotZero(t, result.ResultData.ID)

	// Verify savedgames row exists
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM savedgames WHERE id=$1", result.ResultData.ID)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	// Verify junction rows exist
	err = db.Get(&count, "SELECT COUNT(*) FROM savedgamesplayers WHERE saved_game_id=$1", result.ResultData.ID)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestAddStandardSavedGame_Integration_RollbackOnBadPlayer(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	locationID := getLocationID(t, "Elmwood")
	playerOneID := getPlayerID(t, "playerone")
	badPlayerID := 564523 // does not exist → FK violation

	savedGame := models.SavedGame{
		TotalPoints:       900,
		AveragePoints:     300,
		WinningPlayerName: sql.NullString{String: "playerone"},
		WinningPlayerId:   sql.NullInt32{Int32: int32(playerOneID), Valid: true},
		LocationId:        locationID,
		Players: []models.Player{
			{ID: playerOneID, PlayerName: "playerone", Score: 500},
			{ID: badPlayerID, PlayerName: "ghost", Score: 400},
		},
	}

	// Count before
	var before int
	_ = db.Get(&before, "SELECT COUNT(*) FROM savedgames")

	result := repo.addStandardSavedGame(savedGame)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)

	// Verify NOTHING inserted
	var after int
	_ = db.Get(&after, "SELECT COUNT(*) FROM savedgames")
	require.Equal(t, before, after)
}

func TestAddStandardSavedGame_Integration_RollbackOnBadWinningPlayer(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	locationID := getLocationID(t, "Elmwood")

	savedGame := models.SavedGame{
		TotalPoints:       900,
		AveragePoints:     300,
		WinningPlayerName: sql.NullString{String: "ghost"},
		WinningPlayerId:   sql.NullInt32{Int32: 8765439, Valid: true}, // invalid FK
		LocationId:        locationID,
		Players:           []models.Player{
			{ID: 1, PlayerName: "playerone", Score: 500},
			{ID: 2, PlayerName: "ghost", Score: 400},
		},
	}

	var before int
	_ = db.Get(&before, "SELECT COUNT(*) FROM savedgames")

	result := repo.addStandardSavedGame(savedGame)

	require.NotNil(t, result.Err)
	require.Equal(t, 500, result.StatusCode)

	var after int
	_ = db.Get(&after, "SELECT COUNT(*) FROM savedgames")
	require.Equal(t, before, after)
}

func TestAddStandardSavedGame_Integration_RollbackOnBadLocation(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	playerOneID := getPlayerID(t, "playerone")

	savedGame := models.SavedGame{
		TotalPoints:       900,
		AveragePoints:     300,
		WinningPlayerName: sql.NullString{String: "playerone"},
		WinningPlayerId:   sql.NullInt32{Int32: int32(playerOneID), Valid: true},
		LocationId:        999999, // invalid
		Players:           []models.Player{},
	}

	var before int
	_ = db.Get(&before, "SELECT COUNT(*) FROM savedgames")

	result := repo.addStandardSavedGame(savedGame)

	require.NotNil(t, result.Err)

	var after int
	_ = db.Get(&after, "SELECT COUNT(*) FROM savedgames")
	require.Equal(t, before, after)
}

// Verifies full successful transaction: saved game and junction rows inserted.
func TestAddTeamSavedGame_Integration_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)
	locationID := getLocationID(t, "Elmwood")

	// Elmwood has team with id = 1 from seeding
	savedGame := models.SavedGame{
		TotalPoints:   3000,
		AveragePoints: 750,
		WinningTeamId: sql.NullInt32{Int32: 1, Valid: true},
		LocationId:    locationID,
		Teams: []models.Team{
			{ID: 1, Score: 1000},
			{ID: 2, Score: 900},
		},
	}

	result := repo.addTeamSavedGame(savedGame)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)

	// Verify saved game exists
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM savedgames WHERE id=$1", result.ResultData.ID)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	// Verify junction rows
	err = db.Get(&count,
		"SELECT COUNT(*) FROM savedgamesteams WHERE saved_game_id=$1",
		result.ResultData.ID,
	)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

// Verifies multiple teams are correctly inserted.
func TestAddTeamSavedGame_Integration_MultipleTeams_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	locationID := getLocationID(t, "Lawrence")

	savedGame := models.SavedGame{
		TotalPoints:   5000,
		AveragePoints: 1250,
		WinningTeamId: sql.NullInt32{Int32: 2, Valid: true},
		LocationId:    locationID,
		Teams: []models.Team{
			{ID: 2, Score: 1500},
			{ID: 3, Score: 1200},
			{ID: 4, Score: 1100},
		},
	}

	result := repo.addTeamSavedGame(savedGame)

	require.Nil(t, result.Err)

	var count int
	err := db.Get(&count,
		"SELECT COUNT(*) FROM savedgamesteams WHERE saved_game_id=$1",
		result.ResultData.ID,
	)
	require.NoError(t, err)
	require.Equal(t, 3, count)
}


// Verifies transaction rolls back if winning_team_id violates FK.
func TestAddTeamSavedGame_Integration_RollbackOnBadWinningTeam_Unhappy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)
	locationID := getLocationID(t, "Elmwood")

	var before int
	_ = db.Get(&before, "SELECT COUNT(*) FROM savedgames")

	savedGame := models.SavedGame{
		TotalPoints:   1000,
		AveragePoints: 250,
		WinningTeamId: sql.NullInt32{Int32: 999999, Valid: true}, // invalid FK
		LocationId:    locationID,
		Teams:         []models.Team{},
	}

	result := repo.addTeamSavedGame(savedGame)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)

	var after int
	_ = db.Get(&after, "SELECT COUNT(*) FROM savedgames")
	require.Equal(t, before, after) // rollback occurred
}

// Verifies rollback if a team insert in the junction table fails.
func TestAddTeamSavedGame_Integration_RollbackOnBadTeamInJunction_Unhappy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db, nil)

	locationID := getLocationID(t, "Elmwood")

	var before int
	_ = db.Get(&before, "SELECT COUNT(*) FROM savedgames")

	savedGame := models.SavedGame{
		TotalPoints:   2000,
		AveragePoints: 500,
		WinningTeamId: sql.NullInt32{Int32: 1, Valid: true},
		LocationId:    locationID,
		Teams: []models.Team{
			{ID: 1, Score: 1000},
			{ID: 999999, Score: 500}, // invalid FK here
		},
	}

	result := repo.addTeamSavedGame(savedGame)

	require.NotNil(t, result.Err)

	var after int
	_ = db.Get(&after, "SELECT COUNT(*) FROM savedgames")
	require.Equal(t, before, after) // full rollback
}