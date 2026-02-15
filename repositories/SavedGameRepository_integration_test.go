package repositories

import (
	"JeopardyScoreBoardV2/encryption"
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/utils"
	"database/sql"
	"encoding/base64"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Helper to create encryption service for tests
func getTestEncryptionService(t *testing.T) *encryption.EncryptionService {
	key := os.Getenv("ENCRYPTION_KEY")
	keyB64, err := base64.StdEncoding.DecodeString(key)

	if err != nil {
		t.Fatal(err)
	}
	
	if string(keyB64) == "" {
		t.Skip("ENCRYPTION_KEY not set, skipping encryption tests")
	}

	return encryption.NewService(keyB64)
}

//=======================================
// GET / tests retrieving saved games
//======================================

// GetAllSavedGamesDB_Integration_Happy
// Verifies all seeded saved games are returned from the real test database.
func TestGetAllSavedGamesDB_Integration_Happy(t *testing.T) {
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)
	result := repo.GetAllSavedGamesDB()

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)

	// From migrations: 8 total saved games
	require.GreaterOrEqual(t, len(result.ResultData), 8)
}

// GetAllSavedGamesFromLocationDB_Integration_Happy_Elmwood
// Verifies Elmwood returns the 2 seeded saved games.
func TestGetAllSavedGamesFromLocationDB_Integration_Happy_Elmwood(t *testing.T) {
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

	result := repo.GetAllSavedGamesFromLocationDB("Elmwood")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)

	// Verify at least one has decrypted winning player name
	hasDecryptedName := false
	for _, sg := range result.ResultData {
		if sg.WinningPlayerName != "" {
			hasDecryptedName = true
			break
		}
	}
	require.True(t, hasDecryptedName, "at least one game should have decrypted winning player name")
}

// GetAllSavedGamesFromLocationDB_Integration_Happy_Lawrence
// Verifies Lawrence returns the 4 seeded saved games.
func TestGetAllSavedGamesFromLocationDB_Integration_Happy_Lawrence(t *testing.T) {
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)
	result := repo.GetAllSavedGamesFromLocationDB("Lawrence")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 4)
}

// GetAllSavedGamesFromLocationDB_Integration_Happy_Flushing
// Verifies Flushing returns the 2 seeded saved games.
func TestGetAllSavedGamesFromLocationDB_Integration_Happy_Flushing(t *testing.T) {
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

	result := repo.GetAllSavedGamesFromLocationDB("Flushing")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)
}

// GetAllSavedGamesFromLocationDB_Integration_Unhappy_NoGames
// Verifies empty slice when location exists but has no saved games.
func TestGetAllSavedGamesFromLocationDB_Integration_Unhappy_NoGames(t *testing.T) {
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)
	result := repo.GetAllSavedGamesFromLocationDB("Pelham Bay")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Empty(t, result.ResultData)
}

// GetAllSavedGamesFromLocationDB_Integration_Unhappy_BadLocation
// Verifies empty slice when location name does not exist.
func TestGetAllSavedGamesFromLocationDB_Integration_Unhappy_BadLocation(t *testing.T) {
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

	result := repo.GetAllSavedGamesFromLocationDB("NotARealLocation")

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Empty(t, result.ResultData)
}

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
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

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
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)
	fakeID := "999"
	result := repo.DeleteSavedGameDB(fakeID)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusNotFound, result.StatusCode)
}

//=======================================
// INSERT - tests inserting saved games
//=======================================

func getLocationID(t *testing.T, name string) int {
	var id int
	err := db.Get(&id, "SELECT id FROM locations WHERE location_name=$1", name)
	require.NoError(t, err)
	return id
}

func getPlayerIDByHash(t *testing.T, playerName string) int {
	hash := utils.NameHash(playerName)
	var id int
	err := db.Get(&id, "SELECT id FROM players WHERE player_name_hash=$1", hash)
	require.NoError(t, err)
	return id
}

func TestAddStandardSavedGame_Integration_Happy(t *testing.T) {
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

	locationID := getLocationID(t, "Elmwood")
	
	// Use actual player names from seed data
	winnerName := "goofer boofer"
	playerOneName := "goofer boofer"
	playerTwoName := "player mah"
	
	playerOneID := getPlayerIDByHash(t, playerOneName)
	playerTwoID := getPlayerIDByHash(t, playerTwoName)

	// Encrypt the winning player name
	encryptedWinnerName, err := encService.Encrypt(winnerName)
	require.NoError(t, err)
	winnerHash := utils.NameHash(winnerName)

	savedGame := models.SavedGame{
		TotalPoints:                900,
		AveragePoints:              300,
		WinningPlayerNameEncrypted: encryptedWinnerName,
		WinningPlayerNameHash:      winnerHash,
		WinningPlayerId:            sql.NullInt32{Int32: int32(playerOneID), Valid: true},
		LocationId:                 locationID,
		Players: []models.Player{
			{ID: playerOneID, PlayerName: playerOneName, Score: 500},
			{ID: playerTwoID, PlayerName: playerTwoName, Score: 400},
		},
	}

	result := repo.addStandardSavedGame(savedGame)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.NotZero(t, result.ResultData.ID)

	// Verify savedgames row exists with encrypted data
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM savedgames WHERE id=$1", result.ResultData.ID)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	// Verify the encrypted name and hash were stored
	
	storedGame := models.SavedGame{}

	err = db.Get(&storedGame, "SELECT * FROM savedgames WHERE id=$1", result.ResultData.ID)
	require.NoError(t, err)
	require.NotEmpty(t, storedGame.WinningPlayerNameEncrypted)
	require.NotEmpty(t, storedGame.WinningPlayerNameHash)

	// Verify junction rows exist with encrypted player names
	err = db.Get(&count, "SELECT COUNT(*) FROM savedgamesplayers WHERE saved_game_id=$1", result.ResultData.ID)
	require.NoError(t, err)
	require.Equal(t, 2, count)

	// Verify player names are encrypted in junction table
	var junctionPlayer struct {
		PlayerNameEncrypted []byte `db:"player_name_encrypted"`
		PlayerNameHash      []byte `db:"player_name_hash"`
	}
	err = db.Get(&junctionPlayer,
		"SELECT player_name_encrypted, player_name_hash FROM savedgamesplayers WHERE saved_game_id=$1",
		result.ResultData.ID)
		
	require.NoError(t, err)
	require.NotEmpty(t, junctionPlayer.PlayerNameEncrypted)
	require.NotEmpty(t, junctionPlayer.PlayerNameHash)
}

func TestAddStandardSavedGame_Integration_RollbackOnBadPlayer(t *testing.T) {
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

	locationID := getLocationID(t, "Elmwood")
	playerOneName := "goofer boofer"
	playerOneID := getPlayerIDByHash(t, playerOneName)
	badPlayerID := 564523 // does not exist → FK violation

	winnerName := "goofer boofer"
	encryptedWinnerName, _ := encService.Encrypt(winnerName)
	winnerHash := utils.NameHash(winnerName)

	savedGame := models.SavedGame{
		TotalPoints:                900,
		AveragePoints:              300,
		WinningPlayerNameEncrypted: encryptedWinnerName,
		WinningPlayerNameHash:      winnerHash,
		WinningPlayerId:            sql.NullInt32{Int32: int32(playerOneID), Valid: true},
		LocationId:                 locationID,
		Players: []models.Player{
			{ID: playerOneID, PlayerName: playerOneName, Score: 500},
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
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

	locationID := getLocationID(t, "Elmwood")

	winnerName := "ghost"
	encryptedWinnerName, _ := encService.Encrypt(winnerName)
	winnerHash := utils.NameHash(winnerName)

	savedGame := models.SavedGame{
		TotalPoints:                900,
		AveragePoints:              300,
		WinningPlayerNameEncrypted: encryptedWinnerName,
		WinningPlayerNameHash:      winnerHash,
		WinningPlayerId:            sql.NullInt32{Int32: 8765439, Valid: true}, // invalid FK
		LocationId:                 locationID,
		Players:                    []models.Player{},
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
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

	playerOneName := "goofer boofer"
	playerOneID := getPlayerIDByHash(t, playerOneName)

	winnerName := "goofer boofer"
	encryptedWinnerName, _ := encService.Encrypt(winnerName)
	winnerHash := utils.NameHash(winnerName)

	savedGame := models.SavedGame{
		TotalPoints:                900,
		AveragePoints:              300,
		WinningPlayerNameEncrypted: encryptedWinnerName,
		WinningPlayerNameHash:      winnerHash,
		WinningPlayerId:            sql.NullInt32{Int32: int32(playerOneID), Valid: true},
		LocationId:                 999999, // invalid
		Players:                    []models.Player{},
	}

	var before int
	_ = db.Get(&before, "SELECT COUNT(*) FROM savedgames")

	result := repo.addStandardSavedGame(savedGame)

	require.NotNil(t, result.Err)

	var after int
	_ = db.Get(&after, "SELECT COUNT(*) FROM savedgames")
	require.Equal(t, before, after)
}

func TestAddStandardSavedGame_Integration_EncryptionError(t *testing.T) {
	// Create repo with bad encryption key
	badService := encryption.NewService([]byte("short")) // Too short, will fail
	repo := GetSqlSavedGameRepository(db, badService)

	locationID := getLocationID(t, "Elmwood")
	playerOneName := "goofer boofer"
	playerOneID := getPlayerIDByHash(t, playerOneName)

	// This will fail during encryption
	savedGame := models.SavedGame{
		TotalPoints:     900,
		AveragePoints:   300,
		WinningPlayerId: sql.NullInt32{Int32: int32(playerOneID), Valid: true},
		LocationId:      locationID,
		Players: []models.Player{
			{ID: playerOneID, PlayerName: playerOneName, Score: 500},
		},
	}

	result := repo.addStandardSavedGame(savedGame)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// Verifies full successful transaction: saved game and junction rows inserted.
func TestAddTeamSavedGame_Integration_Happy(t *testing.T) {
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)
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
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

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
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)
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
	encService := getTestEncryptionService(t)
	repo := GetSqlSavedGameRepository(db, encService)

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