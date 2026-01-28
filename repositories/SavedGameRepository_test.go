package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"database/sql"
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




////////////////////////////
// DELETE
/////////////////////////////

// DeleteSavedGameDB_Happy
// Tests successful deletion when the saved game exists and 1 row is affected.
func TestDeleteSavedGameDB_Happy(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)

	savedGameID := "sg-123"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result := repo.DeleteSavedGameDB(savedGameID)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Contains(t, result.ResultData, savedGameID)
	require.NoError(t, mock.ExpectationsWereMet())
}

// DeleteSavedGameDB_Unhappy_DBExecError
// Tests when the database returns an error during Exec.
func TestDeleteSavedGameDB_Unhappy_DBExecError(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)

	savedGameID := "sg-123"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameID).
		WillReturnError(errors.New("db failure"))

	result := repo.DeleteSavedGameDB(savedGameID)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
	require.Empty(t, result.ResultData)
	require.NoError(t, mock.ExpectationsWereMet())
}

// DeleteSavedGameDB_Unhappy_RowsAffectedError
// Tests when Exec succeeds but RowsAffected() returns an error.
func TestDeleteSavedGameDB_Unhappy_RowsAffectedError(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)

	savedGameID := "sg-123"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameID).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

	result := repo.DeleteSavedGameDB(savedGameID)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
	require.Empty(t, result.ResultData)
	require.NoError(t, mock.ExpectationsWereMet())
}

// DeleteSavedGameDB_Unhappy_NoRowsAffected
// Tests when Exec succeeds but no rows are deleted (saved game not found).
func TestDeleteSavedGameDB_Unhappy_NoRowsAffected(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)

	savedGameID := "sg-999"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeleteSavedGame)).
		WithArgs(savedGameID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	result := repo.DeleteSavedGameDB(savedGameID)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusNotFound, result.StatusCode)
	require.Empty(t, result.ResultData)
	require.NoError(t, mock.ExpectationsWereMet())
}


////////////////////////////
// POSt/CREATE saved game
/////////////////////////////

func mockSavedGame() models.SavedGame {
	return models.SavedGame{
		TotalPoints:       1200,
		AveragePoints:     400,
		WinningPlayerName: sql.NullString{String: "Darien"},
		WinningPlayerId:  sql.NullInt32{Int32: 1},
		LocationId:        1,
		Players: []models.Player{
			{ID: 1, PlayerName: "Darien", Score: 500},
			{ID: 2, PlayerName: "Vicky", Score: 400},
		},
	}
}


// addStandardSavedGame_Happy
// Tests full successful transaction: insert saved game, insert players, commit.
func TestAddStandardSavedGame_Happy(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)
	game := mockSavedGame()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WithArgs(
			game.TotalPoints,
			game.AveragePoints,
			game.WinningPlayerName,
			game.WinningPlayerId,
			game.LocationId,
		).
		WillReturnRows(sqlmock.NewRows([]string{"savedgameid"}).AddRow(game.ID))

	for _, p := range game.Players {
		mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
			WithArgs(p.ID, game.ID, p.Score, p.PlayerName).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	mock.ExpectCommit()

	result := repo.addStandardSavedGame(game)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.Equal(t, game.ID, result.ResultData.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

// addStandardSavedGame_Unhappy_BeginTxFails
// Tests failure to start transaction.
func TestAddStandardSavedGame_Unhappy_BeginTxFails(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)
	game := mockSavedGame()

	mock.ExpectBegin().WillReturnError(errors.New("tx fail"))

	result := repo.addStandardSavedGame(game)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// addStandardSavedGame_Unhappy_InsertSavedGameFails
// Tests failure inserting the saved game before players are inserted.
func TestAddStandardSavedGame_Unhappy_InsertSavedGameFails(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)
	game := mockSavedGame()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnError(errors.New("insert fail"))

	mock.ExpectRollback()

	result := repo.addStandardSavedGame(game)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// addStandardSavedGame_Unhappy_PlayerInsertFails
// Tests failure while inserting one of the players into the junction table.
func TestAddStandardSavedGame_Unhappy_PlayerInsertFails(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)
	game := mockSavedGame()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnRows(sqlmock.NewRows([]string{"savedgameid"}).AddRow("sg-1"))

	// First player succeeds
	mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
		WithArgs("p1", "sg-1", 500, "A").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Second player fails
	mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
		WithArgs("p2", "sg-1", 400, "B").
		WillReturnError(errors.New("player insert fail"))

	mock.ExpectRollback()

	result := repo.addStandardSavedGame(game)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// addStandardSavedGame_Unhappy_CommitFails
// Tests failure during transaction commit after all inserts succeed.
func TestAddStandardSavedGame_Unhappy_CommitFails(t *testing.T) {
	mock, repo := setupSavedGameRepo(t)
	game := mockSavedGame()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnRows(sqlmock.NewRows([]string{"savedgameid"}).AddRow("sg-1"))

	for _, p := range game.Players {
		mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
			WithArgs(p.ID, "sg-1", p.Score, p.PlayerName).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	mock.ExpectCommit().WillReturnError(errors.New("commit fail"))

	result := repo.addStandardSavedGame(game)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}
