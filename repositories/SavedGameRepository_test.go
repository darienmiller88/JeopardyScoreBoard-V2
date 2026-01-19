package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSavedGameRepo(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *sqlSavedGameRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := &sqlSavedGameRepository{db: sqlxDB}

	t.Cleanup(func() {
		db.Close()
	})

	return sqlxDB, mock, repo
}



////////////////////////////
//GET
///////////////////////////
// TestGetAllSavedGamesDB_HappyPath_ReturnsMultipleGames verifies that GetAllSavedGamesDB
// successfully retrieves multiple saved games from the database. It mocks a result set with
// two games containing player winners and validates that all fields are correctly mapped.
func TestGetAllSavedGamesDB_HappyPath_ReturnsMultipleGames(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	now := time.Now()
	expectedGames := []models.SavedGame{
		{
			ID:                1,
			CreatedAt:         now,
			UpdatedAt:         now,
			TotalPoints:       16000,
			AveragePoints:     4000.0,
			LocationId:        1,
			WinningPlayerName: sql.NullString{String: "Alice", Valid: true},
			WinningPlayerId:   sql.NullInt32{Int32: 1, Valid: true},
		},
		{
			ID:                2,
			CreatedAt:         now,
			UpdatedAt:         now,
			TotalPoints:       18000,
			AveragePoints:     4500.0,
			LocationId:        2,
			WinningPlayerName: sql.NullString{String: "Bob", Valid: true},
			WinningPlayerId:   sql.NullInt32{Int32: 2, Valid: true},
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "total_score", "average_score",
		"location_id", "winning_player_name", "winning_team_id", "winning_player_id",
	}).
		AddRow(1, now, now, 16000, 4000.0, 1, "Alice", nil, 1).
		AddRow(2, now, now, 18000, 4500.0, 2, "Bob", nil, 2)

	mock.ExpectQuery(`SELECT \* FROM savedgames`).WillReturnRows(rows)

	// Act
	result := repo.GetAllSavedGamesDB()

	// Assert
	assert.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 2)
	assert.Equal(t, expectedGames[0].ID, result.ResultData[0].ID)
	assert.Equal(t, expectedGames[0].WinningPlayerName, result.ResultData[0].WinningPlayerName)
	assert.Equal(t, expectedGames[1].ID, result.ResultData[1].ID)
	assert.Equal(t, expectedGames[1].WinningPlayerName, result.ResultData[1].WinningPlayerName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesDB_HappyPath_ReturnsEmptyList verifies that when no saved games
// exist in the database, the function returns an empty slice (not nil) with a 200 OK status.
// This ensures the API can handle empty result sets gracefully.
func TestGetAllSavedGamesDB_HappyPath_ReturnsEmptyList(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "total_score", "average_score",
		"location_id", "winning_player_name", "winning_team_id", "winning_player_id",
	})

	mock.ExpectQuery(`SELECT \* FROM savedgames`).WillReturnRows(rows)

	// Act
	result := repo.GetAllSavedGamesDB()

	// Assert
	assert.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.NotNil(t, result.ResultData) // Should be empty slice, not nil
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesDB_HappyPath_ReturnsGameWithTeamWinner verifies that the function
// correctly handles games won by teams rather than individual players. It checks that
// WinningTeamId is populated while WinningPlayerId and WinningPlayerName are null.
func TestGetAllSavedGamesDB_HappyPath_ReturnsGameWithTeamWinner(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "total_score", "average_score",
		"location_id", "winning_player_name", "winning_team_id", "winning_player_id",
	}).AddRow(1, now, now, 20000, 5000.0, 1, nil, 5, nil)

	mock.ExpectQuery(`SELECT \* FROM savedgames`).WillReturnRows(rows)

	// Act
	result := repo.GetAllSavedGamesDB()

	// Assert
	assert.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 1)
	assert.Equal(t, int32(5), result.ResultData[0].WinningTeamId.Int32)
	assert.True(t, result.ResultData[0].WinningTeamId.Valid)
	assert.False(t, result.ResultData[0].WinningPlayerId.Valid)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesDB_UnhappyPath_DatabaseError verifies that when the database
// encounters a connection error, the function returns a 500 Internal Server Error
// with an appropriate error message and empty result set.
func TestGetAllSavedGamesDB_UnhappyPath_DatabaseError(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	expectedError := fmt.Errorf("database connection lost")
	mock.ExpectQuery(`SELECT \* FROM savedgames`).WillReturnError(expectedError)

	// Act
	result := repo.GetAllSavedGamesDB()

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.Contains(t, result.Err.Error(), "database connection lost")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesDB_UnhappyPath_QueryTimeout tests the error handling when
// a database query exceeds its timeout limit. Ensures the function returns a 500
// status code with a descriptive timeout error message.
func TestGetAllSavedGamesDB_UnhappyPath_QueryTimeout(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	timeoutError := fmt.Errorf("query timeout exceeded")
	mock.ExpectQuery(`SELECT \* FROM savedgames`).WillReturnError(timeoutError)

	// Act
	result := repo.GetAllSavedGamesDB()

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.Contains(t, result.Err.Error(), "query timeout")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesDB_UnhappyPath_ScanError verifies that when database rows
// contain data that cannot be scanned into the SavedGame struct (type mismatch),
// the function handles the scan error gracefully and returns a 500 error.
func TestGetAllSavedGamesDB_UnhappyPath_ScanError(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	// Create rows with wrong data types to cause scan error
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "total_score", "average_score",
		"location_id", "winning_player_name", "winning_team_id", "winning_player_id",
	}).
		AddRow("invalid_id", time.Now(), time.Now(), 16000, 4000.0, 1, "Alice", nil, 1)

	mock.ExpectQuery(`SELECT \* FROM savedgames`).WillReturnRows(rows)

	// Act
	result := repo.GetAllSavedGamesDB()

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesDB_HappyPath_ReturnsGamesWithNullValues verifies that the function
// properly handles games where no winner has been determined yet (all winner fields are null).
// This is a valid state for games that may still be in progress or incomplete.
func TestGetAllSavedGamesDB_HappyPath_ReturnsGamesWithNullValues(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "total_score", "average_score",
		"location_id", "winning_player_name", "winning_team_id", "winning_player_id",
	}).
		AddRow(1, now, now, 16000, 4000.0, 1, nil, nil, nil)

	mock.ExpectQuery(`SELECT \* FROM savedgames`).WillReturnRows(rows)

	// Act
	result := repo.GetAllSavedGamesDB()

	// Assert
	assert.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 1)
	assert.False(t, result.ResultData[0].WinningTeamId.Valid)
	assert.False(t, result.ResultData[0].WinningPlayerId.Valid)
	assert.False(t, result.ResultData[0].WinningPlayerName.Valid)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesDB_UnhappyPath_ClosedDatabase verifies that when attempting
// to query a closed database connection, the function returns an appropriate error
// and 500 status code rather than panicking.
func TestGetAllSavedGamesDB_UnhappyPath_ClosedDatabase(t *testing.T) {
	// Arrange
	db, mock, repo := setupSavedGameRepo(t)
	db.Close() // Close database before query

	// Act
	result := repo.GetAllSavedGamesDB()

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	_ = mock
}

// TestGetAllSavedGamesFromLocationDB_HappyPath_OneResult verifies that when querying
// for saved games at a specific location, the function correctly retrieves and maps
// a single game with all its associated data including winner information.
func TestGetAllSavedGamesFromLocationDB_HappyPath_OneResult(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"total_score",
		"average_score",
		"location_id",
		"winning_player_name",
		"winning_team_id",
		"winning_player_id",
	}).AddRow(
		1,
		time.Now(),
		time.Now(),
		4200,
		1400.0,
		10,
		sql.NullString{String: "Alice", Valid: true},
		sql.NullInt32{Int32: 2, Valid: true},
		sql.NullInt32{Int32: 5, Valid: true},
	)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs("Elmwood").
		WillReturnRows(rows)

	result := repo.GetAllSavedGamesFromLocationDB("Elmwood")

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 1)
	require.Equal(t, 4200, result.ResultData[0].TotalPoints)
}

// TestGetAllSavedGamesFromLocationDB_HappyPath_MultipleResults verifies that the function
// can retrieve multiple saved games from the same location and correctly maps all rows
// into the result slice.
func TestGetAllSavedGamesFromLocationDB_HappyPath_MultipleResults(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"total_score",
		"average_score",
		"location_id",
		"winning_player_name",
		"winning_team_id",
		"winning_player_id",
	}).AddRow(1, time.Now(), time.Now(), 3000, 1000, 5, nil, nil, nil).
		AddRow(2, time.Now(), time.Now(), 4500, 1500, 5, nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs("Lawrence").
		WillReturnRows(rows)

	result := repo.GetAllSavedGamesFromLocationDB("Lawrence")

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)
}

// TestGetAllSavedGamesFromLocationDB_HappyPath_NoResults verifies that when querying
// for a location with no saved games, the function returns an empty slice (not nil)
// with a successful 200 OK status, allowing the API to handle empty results gracefully.
func TestGetAllSavedGamesFromLocationDB_HappyPath_NoResults(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"total_score",
		"average_score",
		"location_id",
		"winning_player_name",
		"winning_team_id",
		"winning_player_id",
	})

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs("Empty Location").
		WillReturnRows(rows)

	result := repo.GetAllSavedGamesFromLocationDB("Empty Location")

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 0)
}

// TestGetAllSavedGamesFromLocationDB_HappyPath_MixedPlayerAndTeamGames verifies that
// when a location has games won by both individual players and teams, the function
// correctly differentiates between them based on which winner fields are populated.
func TestGetAllSavedGamesFromLocationDB_HappyPath_MixedPlayerAndTeamGames(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	locationName := "Elmwood"
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "total_score", "average_score",
		"location_id", "winning_player_name", "winning_team_id", "winning_player_id",
	}).
		AddRow(1, now, now, 16000, 4000.0, 1, "Alice", nil, 1).
		AddRow(2, now, now, 20000, 5000.0, 1, nil, 5, nil).
		AddRow(3, now, now, 18000, 4500.0, 1, "Bob", nil, 2)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs(locationName).
		WillReturnRows(rows)

	// Act
	result := repo.GetAllSavedGamesFromLocationDB(locationName)

	// Assert
	assert.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 3)
	
	// First game: player winner
	assert.True(t, result.ResultData[0].WinningPlayerName.Valid)
	assert.Equal(t, "Alice", result.ResultData[0].WinningPlayerName.String)
	assert.False(t, result.ResultData[0].WinningTeamId.Valid)
	
	// Second game: team winner
	assert.False(t, result.ResultData[1].WinningPlayerName.Valid)
	assert.True(t, result.ResultData[1].WinningTeamId.Valid)
	assert.Equal(t, int32(5), result.ResultData[1].WinningTeamId.Int32)
	
	// Third game: player winner
	assert.True(t, result.ResultData[2].WinningPlayerName.Valid)
	assert.Equal(t, "Bob", result.ResultData[2].WinningPlayerName.String)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesFromLocationDB_HappyPath_NullWinnerFields verifies that games
// with no determined winner (all winner fields null) are handled correctly. This might
// represent games that are incomplete or where the winner hasn't been recorded yet.
func TestGetAllSavedGamesFromLocationDB_HappyPath_NullWinnerFields(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"total_score",
		"average_score",
		"location_id",
		"winning_player_name",
		"winning_team_id",
		"winning_player_id",
	}).AddRow(
		1,
		time.Now(),
		time.Now(),
		2000,
		1000,
		3,
		sql.NullString{Valid: false},
		sql.NullInt32{Valid: false},
		sql.NullInt32{Valid: false},
	)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs("No Winners").
		WillReturnRows(rows)

	result := repo.GetAllSavedGamesFromLocationDB("No Winners")

	require.NoError(t, result.Err)
	require.False(t, result.ResultData[0].WinningPlayerName.Valid)
	require.False(t, result.ResultData[0].WinningTeamId.Valid)
}

// TestGetAllSavedGamesFromLocationDB_QueryError verifies that when the database
// query fails (e.g., connection issues, syntax errors), the function returns a
// 500 Internal Server Error with an empty result set.
func TestGetAllSavedGamesFromLocationDB_UnhappyPath_QueryError(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs("Elmwood").
		WillReturnError(errors.New("db failure"))

	result := repo.GetAllSavedGamesFromLocationDB("Elmwood")

	require.Error(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
	require.Len(t, result.ResultData, 0)
}

// TestGetAllSavedGamesFromLocationDB_SubqueryMultipleRows tests the error handling
// when a subquery in the SQL statement unexpectedly returns multiple rows. This
// typically indicates a database integrity issue and should be handled as an error.
func TestGetAllSavedGamesFromLocationDB_UnhappyPath_SubqueryMultipleRows(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs("Duplicate Location").
		WillReturnError(errors.New("more than one row returned by a subquery"))

	result := repo.GetAllSavedGamesFromLocationDB("Duplicate Location")

	require.Error(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestGetAllSavedGamesFromLocationDB_UnhappyPath_ScanError verifies that when
// database rows contain data types that don't match the expected struct fields
// (e.g., string where int is expected), the function handles the scan error
// gracefully and returns a 500 error.
func TestGetAllSavedGamesFromLocationDB_UnhappyPath_ScanError(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	locationName := "Elmwood"

	// Create rows with wrong data type to cause scan error
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "total_score", "average_score",
		"location_id", "winning_player_name", "winning_team_id", "winning_player_id",
	}).
		AddRow("invalid_id", time.Now(), time.Now(), 16000, 4000.0, 1, "Alice", nil, 1)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs(locationName).
		WillReturnRows(rows)

	// Act
	result := repo.GetAllSavedGamesFromLocationDB(locationName)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesFromLocationDB_EmptyLocationName verifies that when an empty
// string is provided as the location name, the query executes but returns no results.
// This is a valid edge case that should be handled gracefully without errors.
func TestGetAllSavedGamesFromLocationDB_EmptyLocationName(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)
	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"total_score",
		"average_score",
		"location_id",
		"winning_player_name",
		"winning_team_id",
		"winning_player_id",
	})

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs("").
		WillReturnRows(rows)

	result := repo.GetAllSavedGamesFromLocationDB("")

	require.NoError(t, result.Err)
	require.Len(t, result.ResultData, 0)
}









////////////////////////////
// DELETE 
/////////////////////////////

// TestDeleteSavedGameDB_HappyPath_DeletesExistingGame verifies that when a valid
// saved game ID is provided, the function successfully deletes the game from the
// database and returns a 200 OK status with a success message confirming the deletion.
func TestDeleteSavedGameDB_HappyPath_DeletesExistingGame(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := "1"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Contains(t, result.ResultData, "Saved game 1 deleted successfully")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_HappyPath_DeletesGameWithHighId verifies that the function
// correctly handles deletion of saved games with large ID values (edge case for ID
// range) and returns appropriate success messages containing the specific ID.
func TestDeleteSavedGameDB_HappyPath_DeletesGameWithHighId(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := "99999"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Contains(t, result.ResultData, "99999")
	assert.Contains(t, result.ResultData, "deleted successfully")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_HappyPath_DeletesGameWithSpecialCaseId tests the edge case
// of deleting a game with ID of 0, which is a boundary value that might be handled
// specially in some systems. Verifies that ID 0 is treated as a valid ID.
func TestDeleteSavedGameDB_HappyPath_DeletesGameWithSpecialCaseId(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := "0" // Edge case: ID of 0

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Contains(t, result.ResultData, "Saved game 0 deleted successfully")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_UnhappyPath_DatabaseError verifies that when the database
// encounters a connection error during deletion, the function returns a 500 Internal
// Server Error with the appropriate error message and no result data.
func TestDeleteSavedGameDB_UnhappyPath_DatabaseError(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := "1"
	expectedError := fmt.Errorf("database connection lost")

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnError(expectedError)

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.Contains(t, result.Err.Error(), "database connection lost")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_UnhappyPath_InvalidIdFormat verifies that when a non-numeric
// ID is provided (which PostgreSQL cannot cast to integer), the function handles the
// database type conversion error gracefully and returns a 500 error.
func TestDeleteSavedGameDB_UnhappyPath_InvalidIdFormat(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := "invalid_id"

	// PostgreSQL would reject this, simulating that error
	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnError(fmt.Errorf("pq: invalid input syntax for integer: \"invalid_id\""))

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.Contains(t, result.Err.Error(), "invalid input syntax")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_UnhappyPath_EmptyId verifies that when an empty string is
// provided as the game ID, the database rejects it with a type conversion error,
// and the function properly handles this error case.
func TestDeleteSavedGameDB_UnhappyPath_EmptyId(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := ""

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnError(fmt.Errorf("pq: invalid input syntax for integer: \"\""))

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_UnhappyPath_NegativeId verifies that when attempting to delete
// a game with a negative ID (which likely doesn't exist), the function detects that
// no rows were affected and returns a 404 Not Found error with an appropriate message.
func TestDeleteSavedGameDB_UnhappyPath_NegativeId(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := "-1"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.Contains(t, result.Err.Error(), "saved game not found by id: -1")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_UnhappyPath_RowsAffectedError verifies that when the database
// cannot determine how many rows were affected by the delete operation (internal DB error),
// the function handles this gracefully and returns a 500 error.
func TestDeleteSavedGameDB_UnhappyPath_RowsAffectedError(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := "1"
	rowsAffectedError := fmt.Errorf("failed to get rows affected")

	// Return a result that will error when RowsAffected() is called
	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnResult(sqlmock.NewErrorResult(rowsAffectedError))

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.Contains(t, result.Err.Error(), "failed to get rows affected")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_UnhappyPath_ForeignKeyViolation verifies error handling when
// deletion would violate foreign key constraints. Although CASCADE should typically
// prevent this, this test ensures the function handles constraint violations gracefully
// in case of misconfigured database schemas or missing CASCADE directives.
func TestDeleteSavedGameDB_UnhappyPath_ForeignKeyViolation(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	savedGameId := "1"

	// Simulate foreign key constraint violation (though CASCADE should prevent this)
	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnError(fmt.Errorf("pq: foreign key constraint violation"))

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.Contains(t, result.Err.Error(), "foreign key constraint")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_UnhappyPath_ConnectionTimeout verifies that when the database
// connection times out during the delete operation, the function handles the timeout
// error properly and returns a 500 status with a descriptive error message.
func TestDeleteSavedGameDB_UnhappyPath_ConnectionTimeout(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := "1"
	timeoutError := fmt.Errorf("connection timeout")

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnError(timeoutError)

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.Contains(t, result.Err.Error(), "connection timeout")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_UnhappyPath_ClosedDatabase verifies that attempting to delete
// a game after the database connection has been closed is handled gracefully. The
// function should return an error rather than panicking.
func TestDeleteSavedGameDB_UnhappyPath_ClosedDatabase(t *testing.T) {
	// Arrange
	db, mock, repo := setupSavedGameRepo(t)
	savedGameId := "1"

	db.Close() // Close database before deletion

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	_ = mock
}

// TestDeleteSavedGameDB_UnhappyPath_SqlInjectionAttempt verifies that SQL injection
// attempts (e.g., "1 OR 1=1") are safely handled by parameterized queries. The database
// should reject the malformed input as an invalid integer, preventing the injection attack.
func TestDeleteSavedGameDB_UnhappyPath_SqlInjectionAttempt(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)

	// SQL injection attempt
	savedGameId := "1 OR 1=1"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnError(fmt.Errorf("pq: invalid input syntax for integer"))

	// Act
	result := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteSavedGameDB_UnhappyPath_MultipleCallsSameId verifies idempotency behavior
// when attempting to delete the same game twice. The first deletion should succeed,
// but the second attempt should fail with a 404 Not Found since the game no longer exists.
// This tests the proper handling of "already deleted" scenarios.
func TestDeleteSavedGameDB_UnhappyPath_MultipleCallsSameId(t *testing.T) {
	// Arrange
	_, mock, repo := setupSavedGameRepo(t)
	savedGameId := "1"

	// First delete succeeds
	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Second delete fails (game already deleted)
	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameId).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Act
	result1 := repo.DeleteSavedGameDB(savedGameId)
	result2 := repo.DeleteSavedGameDB(savedGameId)

	// Assert
	// First delete succeeds
	assert.NoError(t, result1.Err)
	assert.Equal(t, http.StatusOK, result1.StatusCode)

	// Second delete fails
	assert.Error(t, result2.Err)
	assert.Equal(t, http.StatusNotFound, result2.StatusCode)
	assert.Contains(t, result2.Err.Error(), "saved game not found")

	assert.NoError(t, mock.ExpectationsWereMet())
}



///////////////////////////
//POST/ADD
//////////////////////////


// Tests:
// - WinningPlayerId is valid (player-based game)
// - Saved game insert succeeds
// - One player inserted into savedgameplayers
// - Transaction commits successfully
// - Returns 201 Created with populated SavedGame ID
func TestAddSavedGameDB_HappyPath_PlayerWin_SinglePlayer(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)
	savedGame := models.SavedGame{
		TotalPoints:    3000,
		AveragePoints: 1500,
		WinningPlayerId: sql.NullInt32{Int32: 1, Valid: true},
		WinningPlayerName: sql.NullString{String: "Alice", Valid: true},
		LocationId:     2,
		Players: []models.Player{
			{ID: 1, Score: 1500},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WithArgs(
			savedGame.TotalPoints,
			savedGame.AveragePoints,
			savedGame.WinningPlayerName,
			savedGame.WinningPlayerName,
			savedGame.LocationId,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
		WithArgs(1, 42, 1500, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	result := repo.AddSavedGameDB(savedGame)

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.Equal(t, 42, result.ResultData.ID)
}

// Tests:
// - Multiple players added to savedgameplayers
// - Loop executes fully
// - Transaction commits once all inserts succeed
func TestAddSavedGameDB_HappyPath_PlayerWin_MultiplePlayers(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)
	savedGame := models.SavedGame{
		TotalPoints:    6000,
		AveragePoints: 2000,
		WinningPlayerId: sql.NullInt32{Int32: 2, Valid: true},
		WinningPlayerName: sql.NullString{String: "Bob", Valid: true},
		LocationId:     1,
		Players: []models.Player{
			{ID: 1, Score: 2000},
			{ID: 2, Score: 2500},
			{ID: 3, Score: 1500},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(99))

	for _, p := range savedGame.Players {
		mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
			WithArgs(p.ID, 99, p.Score, p.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	result := repo.AddSavedGameDB(savedGame)

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.Equal(t, 99, result.ResultData.ID)
}

////////////
//UNHAPPY PATHS - Inserting Player saved game
////////////

// Tests:
// - Failure to begin transaction
// - Immediate 500 response
// - No queries executed
func TestAddSavedGameDB_UnhappyPath_PlayerWin_BeginTransactionFails(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New("tx begin failed"))

	result := repo.AddSavedGameDB(models.SavedGame{
		WinningPlayerId: sql.NullInt32{Int32: 1, Valid: true},
	})

	require.Error(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// Tests:
// - FK violation on savedgames insert (missing location or player)
// - Postgres error code 23503
// - Returns 404 Not Found
// - Transaction rolled back
func TestAddSavedGameDB_UnhappyPath_PlayerWin_InsertSavedGameFKViolation(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnError(&pq.Error{Code: "23503"})

	mock.ExpectRollback()

	result := repo.AddSavedGameDB(models.SavedGame{
		WinningPlayerId:   sql.NullInt32{Int32: 1, Valid: true},
		WinningPlayerName: sql.NullString{String: "Ghost", Valid: true},
		LocationId:        999,
	})

	require.Error(t, result.Err)
	require.Equal(t, http.StatusNotFound, result.StatusCode)
}

// Tests:
// - Saved game insert succeeds
// - Player insert fails due to missing player (FK violation)
// - Returns 404
// - Transaction rolled back
func TestAddSavedGameDB_UnhappyPath_PlayerWin_PlayerInsertFKViolation(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
		WillReturnError(&pq.Error{Code: "23503"})

	mock.ExpectRollback()

	result := repo.AddSavedGameDB(models.SavedGame{
		WinningPlayerId:   sql.NullInt32{Int32: 1, Valid: true},
		WinningPlayerName: sql.NullString{String: "Alice", Valid: true},
		LocationId:        1,
		Players: []models.Player{
			{ID: 999, Score: 1000},
		},
	})

	require.Error(t, result.Err)
	require.Equal(t, http.StatusNotFound, result.StatusCode)
}

// Tests:
// - Duplicate (player_id, saved_game_id) entry
// - Unique constraint violation (23505)
// - Returns 409 Conflict
// - Transaction rolled back
func TestAddSavedGameDB_UnhappyPath_PlayerWin_DuplicatePlayerGameEntry(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(50))

	mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
		WillReturnError(&pq.Error{Code: "23505"})

	mock.ExpectRollback()

	result := repo.AddSavedGameDB(models.SavedGame{
		WinningPlayerId:   sql.NullInt32{Int32: 1, Valid: true},
		WinningPlayerName: sql.NullString{String: "Alice", Valid: true},
		LocationId:        1,
		Players: []models.Player{
			{ID: 1, Score: 1000},
		},
	})

	require.Error(t, result.Err)
	require.Equal(t, http.StatusConflict, result.StatusCode)
	require.Nil(t, mock.ExpectationsWereMet())
}

// Tests:
// - All inserts succeed
// - Commit fails
// - Returns 500
// - Data is not confirmed persisted
func TestAddSavedGameDB_UnhappyPath_PlayerWin_CommitFails(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(77))

	mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

	result := repo.AddSavedGameDB(models.SavedGame{
		WinningPlayerId:   sql.NullInt32{Int32: 1, Valid: true},
		WinningPlayerName: sql.NullString{String: "Alice", Valid: true},
		LocationId:        1,
		Players: []models.Player{
			{ID: 1, Score: 1000},
		},
	})

	require.Error(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

////////////////////////////////////////////
//HAPPY PATHS - inserting team saved game
//////////////////////////////////////////

// Tests:
// - WinningPlayerId is NOT valid (team-based game)
// - Saved game insert succeeds
// - One team inserted into savedgameteams
// - Transaction commits successfully
// - Returns 201 Created with populated SavedGame ID
func TestAddSavedGameDB_HappyPath_TeamWin_SingleTeam(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)
	savedGame := models.SavedGame{
		TotalPoints:   4000,
		AveragePoints: 2000,
		WinningTeamId: sql.NullInt32{Int32: 3, Valid: true},
		LocationId:    4,
		Teams: []models.Team{
			{ID: 3, Score: 4000},
		},
	}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewTeamSavedGame)).
		WithArgs(
			savedGame.TotalPoints,
			savedGame.AveragePoints,
			savedGame.WinningTeamId,
			savedGame.LocationId,
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(101))

	mock.ExpectExec(regexp.QuoteMeta(constants.InsertTeamsForSavedGame)).
		WithArgs(3, 101, 4000).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	result := repo.AddSavedGameDB(savedGame)

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.Equal(t, 101, result.ResultData.ID)
}

// Tests:
// - Multiple teams added to savedgameteams
// - Loop executes fully
// - Transaction commits after all inserts succeed
func TestAddSavedGameDB_HappyPath_TeamWin_MultipleTeams(t *testing.T) {
	_, mock, repo := setupSavedGameRepo(t)
	savedGame := models.SavedGame{
		TotalPoints:    9000,
		AveragePoints: 3001.65,
		WinningTeamId:   sql.NullInt32{Int32: 2, Valid: true},
		LocationId:     6,
		Teams: []models.Team{
			{ID: 1, Score: 3000},
			{ID: 2, Score: 4000},
			{ID: 3, Score: 2000},
		},
	}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewTeamSavedGame)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(202))

	for _, team := range savedGame.Teams {
		mock.ExpectExec(regexp.QuoteMeta(constants.InsertTeamsForSavedGame)).
			WithArgs(team.ID, 202, team.Score).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	result := repo.AddSavedGameDB(savedGame)

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.Equal(t, 202, result.ResultData.ID)
}
