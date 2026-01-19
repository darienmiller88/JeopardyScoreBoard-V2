package repositories

import (
	"JeopardyScoreBoardV2/models"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

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