package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
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

////////////////////////////
//
/////////////////////////////