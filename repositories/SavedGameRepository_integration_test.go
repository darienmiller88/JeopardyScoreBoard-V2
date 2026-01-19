package repositories

import (
	"JeopardyScoreBoardV2/models"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// Test helper function to create a test saved game
func createTestSavedGame(t *testing.T, db *sqlx.DB, id string) {
	// Create a saved game with known data
	query := `
		INSERT INTO savedgames (id, location_id, total_score, average_score)
		VALUES ($1, (SELECT id FROM locations WHERE location_name=$2), $3, $4)
	`
	_, err := db.Exec(query, id, "Elmwood", 4000, 2000.45)
	if err != nil {
		t.Fatalf("Failed to create test saved game: %v", err)
	}
}

// Test helper function to clean up a test saved game
func deleteTestSavedGame(t *testing.T, db *sqlx.DB, savedGameId string) {
	query := `DELETE FROM savedgames WHERE id = $1`
	_, err := db.Exec(query, savedGameId)
	if err != nil {
		t.Logf("Warning: Failed to clean up test saved game %s: %v", savedGameId, err)
	}
}

// Test helper function to verify a saved game exists
func savedGameExists(t *testing.T, db *sqlx.DB, savedGameId string) bool {
	var count int
	query := `SELECT COUNT(*) FROM savedgames WHERE id = $1`
	err := db.Get(&count, query, savedGameId)
	if err != nil {
		t.Fatalf("Failed to check if saved game exists: %v", err)
	}
	return count > 0
}

func TestGetAllSavedGamesDB_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db)

	result := repo.GetAllSavedGamesDB()

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Nil(t, result.Err)
	assert.NotNil(t, result.ResultData)
	assert.IsType(t, []models.SavedGame{}, result.ResultData)
}

func TestGetAllSavedGamesDB_Unhappy(t *testing.T) {
	// Create a repository with a closed/invalid database connection
	closedDB, _ := sqlx.Open("postgres", "invalid_connection_string")
	closedDB.Close()
	repo := GetSqlSavedGameRepository(closedDB)

	result := repo.GetAllSavedGamesDB()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.NotNil(t, result.Err)
	assert.Equal(t, []models.SavedGame{}, result.ResultData)
}

func TestGetAllSavedGamesFromLocationDB_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db)
	locationName := "Elmwood" // Using one of the 8 fixed locations

	result := repo.GetAllSavedGamesFromLocationDB(locationName)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Nil(t, result.Err)
	assert.NotNil(t, result.ResultData)
	assert.IsType(t, []models.SavedGame{}, result.ResultData)
}

func TestGetAllSavedGamesFromLocationDB_Unhappy(t *testing.T) {
	// Create a repository with a closed/invalid database connection
	closedDB, _ := sqlx.Open("postgres", "invalid_connection_string")
	closedDB.Close()
	repo := GetSqlSavedGameRepository(closedDB)

	result := repo.GetAllSavedGamesFromLocationDB("Elmwood")

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.NotNil(t, result.Err)
	assert.Equal(t, []models.SavedGame{}, result.ResultData)
}

func TestGetAllSavedGamesFromLocationDB_EmptyLocation_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db)
	locationName := "" // Empty location name

	result := repo.GetAllSavedGamesFromLocationDB(locationName)

	// This should still return OK with an empty array (no games at empty location)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Nil(t, result.Err)
	assert.NotNil(t, result.ResultData)
}

func TestGetAllSavedGamesFromLocationDB_NonExistentLocation_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db)
	locationName := "NonExistentLocation123456789"

	result := repo.GetAllSavedGamesFromLocationDB(locationName)

	// Should return OK with empty array (no games found at this location)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Nil(t, result.Err)
	assert.NotNil(t, result.ResultData)
	assert.Equal(t, 0, len(result.ResultData))
}

func TestDeleteSavedGameDB_Happy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db)
	id := "1"

	// Create a test saved game
	createTestSavedGame(t, db, id)
	defer deleteTestSavedGame(t, db, id) // Cleanup in case test fails

	// Delete the saved game
	result := repo.DeleteSavedGameDB(id)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Nil(t, result.Err)
	assert.Contains(t, result.ResultData, "deleted successfully")
	assert.Contains(t, result.ResultData, id)

	// Verify it no longer exists
	assert.False(t, savedGameExists(t, db, id))
}

func TestDeleteSavedGameDB_NotFound_Unhappy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db)
	nonExistentId := "456"

	// Verify it doesn't exist
	assert.False(t, savedGameExists(t, db, nonExistentId))

	result := repo.DeleteSavedGameDB(nonExistentId)

	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.NotNil(t, result.Err)
	assert.Contains(t, result.Err.Error(), "saved game not found by id")
	assert.Contains(t, result.Err.Error(), nonExistentId)
	assert.Equal(t, "", result.ResultData)
}

func TestDeleteSavedGameDB_DatabaseError_Unhappy(t *testing.T) {
	closedDB, _ := sqlx.Open("postgres", "invalid_connection_string")
	closedDB.Close()
	repo := GetSqlSavedGameRepository(closedDB)

	result := repo.DeleteSavedGameDB("any-id")

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.NotNil(t, result.Err)
	assert.Equal(t, "", result.ResultData)
}

func TestDeleteSavedGameDB_EmptyId_Unhappy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db)

	result := repo.DeleteSavedGameDB("")

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.NotNil(t, result.Err)
	assert.Equal(t, "", result.ResultData)
}

func TestDeleteSavedGameDB_AlreadyDeleted_Unhappy(t *testing.T) {
	repo := GetSqlSavedGameRepository(db)
	id := "1"

	// Create a test saved game
	createTestSavedGame(t, db, id)
	defer deleteTestSavedGame(t, db, id) // Cleanup

	// Delete once
	result1 := repo.DeleteSavedGameDB(id)
	assert.Equal(t, http.StatusOK, result1.StatusCode)
	assert.False(t, savedGameExists(t, db, id))

	// Try to delete again - should fail with NotFound
	result2 := repo.DeleteSavedGameDB(id)

	assert.Equal(t, http.StatusNotFound, result2.StatusCode)
	assert.NotNil(t, result2.Err)
	assert.Contains(t, result2.Err.Error(), "saved game not found by id")
	assert.Equal(t, "", result2.ResultData)
}
