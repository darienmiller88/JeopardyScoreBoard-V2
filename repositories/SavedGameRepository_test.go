package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"errors"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var savedGameColumns = []string{
	"id",
	"created_at",
	"updated_at",
	"location_id",
	"winning_team_id",
	"winning_player_id",
	"winning_player_name",
	"total_score",
	"average_score",
}

func setupSavedGameRepo(t *testing.T) (sqlmock.Sqlmock, *sqlSavedGameRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := &sqlSavedGameRepository{db: sqlxDB}

	t.Cleanup(func() {
		db.Close()
	})

	return mock, repo
}

func mockSavedGameRows() *sqlmock.Rows {
	return sqlmock.NewRows(savedGameColumns).
		AddRow(
			1,
			time.Now(),
			time.Now(),
			3,
			1,
			5,
			"Darien",
			1200,
			400.0,
		).
		AddRow(
			2,
			time.Now(),
			time.Now(),
			2,
			nil,
			nil,
			"",
			800,
			266.6,
		)
}

////////////////////////////
//GET
///////////////////////////

func TestGetAllSavedGamesDB_Success_Happy(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGames)).
		WillReturnRows(mockSavedGameRows())

	result := repo.GetAllSavedGamesDB()

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)
	require.NoError(t, mock.ExpectationsWereMet())
}


func TestGetAllSavedGamesDB_NoRows_Happy(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGames)).
		WillReturnRows(sqlmock.NewRows(savedGameColumns))

	result := repo.GetAllSavedGamesDB()

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Empty(t, result.ResultData)
}

func TestGetAllSavedGamesFromLocationDB_Happy(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)
	location := "Elmwood"

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs(location).
		WillReturnRows(mockSavedGameRows())

	result := repo.GetAllSavedGamesFromLocationDB(location)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)
}


func TestGetAllSavedGamesDB_Error_Unhappy(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGames)).
		WillReturnError(errors.New("db exploded"))

	result := repo.GetAllSavedGamesDB()

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
	require.Empty(t, result.ResultData)
}

func TestGetAllSavedGamesFromLocationDB_Error_Unhappy(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)
	location := "Flushing"

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs(location).
		WillReturnError(errors.New("db exploded"))

	result := repo.GetAllSavedGamesFromLocationDB(location)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestGetAllSavedGamesFromLocationDB_NoRows_Happy(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)
	location := "Nowhere"

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs(location).
		WillReturnRows(sqlmock.NewRows(savedGameColumns))

	result := repo.GetAllSavedGamesFromLocationDB(location)

	require.Nil(t, result.Err)
	require.Empty(t, result.ResultData)
}


////////////////////////////
// DELETE
/////////////////////////////

// TestDeleteSavedGameDB_HappyPath_DeletesExistingGame verifies that when a valid
// saved game ID is provided, the function successfully deletes the game from the
// database and returns a 200 OK status with a success message confirming the deletion.
