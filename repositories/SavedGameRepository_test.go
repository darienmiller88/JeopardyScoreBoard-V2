package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/encryption"
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

// savedGameColumns defines the columns returned when querying saved games
var savedGameColumns = []string{
	"id",
	"created_at",
	"updated_at",
	"location_id",
	"winning_team_id",
	"winning_player_id",
	"winning_player_name_encrypted",
	"winning_player_name_hash",
	"total_score",
	"average_score",
}

// setupSavedGameRepo creates a mock database and repository for testing
func setupSavedGameRepo(t *testing.T) (sqlmock.Sqlmock, *sqlSavedGameRepository, *encryption.EncryptionService) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	
	// Create test encryption service with 32-byte key
	testKey := []byte("12345678901234567890123456789012")
	encryptionService := encryption.NewService(testKey)
	
	repo := &sqlSavedGameRepository{
		db:                sqlxDB,
		encryptionService: encryptionService,
	}

	t.Cleanup(func() {
		db.Close()
	})

	return mock, repo, encryptionService
}

// mockSavedGameRows returns mock rows with encrypted player names for testing
func mockSavedGameRows(encService *encryption.EncryptionService) *sqlmock.Rows {
	encryptedName1, _ := encService.Encrypt("Darien")
	encryptedName2, _ := encService.Encrypt("")
	
	return sqlmock.NewRows(savedGameColumns).
		AddRow(
			1,
			time.Now(),
			time.Now(),
			3,
			1,
			5,
			encryptedName1,
			[]byte("hash1"),
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
			encryptedName2,
			[]byte{},
			800,
			266.6,
		)
}

////////////////////////////
//GET
///////////////////////////

// TestGetAllSavedGamesDB_Success_Happy verifies that all saved games are returned
// and winning player names are decrypted successfully
func TestGetAllSavedGamesDB_Success_Happy(t *testing.T) {
	mock, repo, encService := setupSavedGameRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGames)).
		WillReturnRows(mockSavedGameRows(encService))

	result := repo.GetAllSavedGamesDB()

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)
	// Verify first game has decrypted winning player name
	require.Equal(t, "Darien", result.ResultData[0].WinningPlayerName)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllSavedGamesDB_NoRows_Happy verifies empty result when no saved games exist
func TestGetAllSavedGamesDB_NoRows_Happy(t *testing.T) {
	mock, repo, _ := setupSavedGameRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGames)).
		WillReturnRows(sqlmock.NewRows(savedGameColumns))

	result := repo.GetAllSavedGamesDB()

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Empty(t, result.ResultData)
}

// TestGetAllSavedGamesFromLocationDB_Happy verifies saved games are returned
// for a specific location with decrypted winning player names
func TestGetAllSavedGamesFromLocationDB_Happy(t *testing.T) {
	mock, repo, encService := setupSavedGameRepo(t)
	location := "Elmwood"

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs(location).
		WillReturnRows(mockSavedGameRows(encService))

	result := repo.GetAllSavedGamesFromLocationDB(location)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)
	require.Equal(t, "Darien", result.ResultData[0].WinningPlayerName)
}

// TestGetAllSavedGamesFromLocationDB_NoRows_Happy verifies empty result when
// location has no saved games
func TestGetAllSavedGamesFromLocationDB_NoRows_Happy(t *testing.T) {
	mock, repo, _ := setupSavedGameRepo(t)
	location := "Nowhere"

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs(location).
		WillReturnRows(sqlmock.NewRows(savedGameColumns))

	result := repo.GetAllSavedGamesFromLocationDB(location)

	require.Nil(t, result.Err)
	require.Empty(t, result.ResultData)
}

// TestGetAllSavedGamesDB_Error_Unhappy verifies proper error handling when
// database query fails
func TestGetAllSavedGamesDB_Error_Unhappy(t *testing.T) {
	mock, repo, _ := setupSavedGameRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGames)).
		WillReturnError(errors.New("db exploded"))

	result := repo.GetAllSavedGamesDB()

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
	require.Empty(t, result.ResultData)
}

// TestGetAllSavedGamesFromLocationDB_Error_Unhappy verifies proper error handling
// when database query fails for location-specific queries
func TestGetAllSavedGamesFromLocationDB_Error_Unhappy(t *testing.T) {
	mock, repo, _ := setupSavedGameRepo(t)
	location := "Flushing"

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGamesFromLocation)).
		WithArgs(location).
		WillReturnError(errors.New("db exploded"))

	result := repo.GetAllSavedGamesFromLocationDB(location)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestGetAllSavedGamesDB_DecryptionError_Unhappy verifies error handling when
// decryption fails on retrieved saved games
func TestGetAllSavedGamesDB_DecryptionError_Unhappy(t *testing.T) {
	mock, repo, _ := setupSavedGameRepo(t)
	
	// Return invalid encrypted data that can't be decrypted
	rows := sqlmock.NewRows(savedGameColumns).
		AddRow(
			1,
			time.Now(),
			time.Now(),
			3,
			nil,
			1,
			[]byte("invalid encrypted data"), // Bad encryption
			[]byte("hash1"),
			1200,
			400.0,
		)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllSavedGames)).
		WillReturnRows(rows)

	result := repo.GetAllSavedGamesDB()

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

////////////////////////////
// DELETE
/////////////////////////////

// TestDeleteSavedGameDB_Happy verifies successful deletion when the saved game
// exists and one row is affected
func TestDeleteSavedGameDB_Happy(t *testing.T) {
	mock, repo, _ := setupSavedGameRepo(t)

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

// TestDeleteSavedGameDB_Unhappy_DBExecError verifies error handling when the
// database returns an error during Exec
func TestDeleteSavedGameDB_Unhappy_DBExecError(t *testing.T) {
	mock, repo, _ := setupSavedGameRepo(t)

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

// TestDeleteSavedGameDB_Unhappy_RowsAffectedError verifies error handling when
// Exec succeeds but RowsAffected() returns an error
func TestDeleteSavedGameDB_Unhappy_RowsAffectedError(t *testing.T) {
	mock, repo, _ := setupSavedGameRepo(t)

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

// TestDeleteSavedGameDB_Unhappy_NoRowsAffected verifies 404 response when Exec
// succeeds but no rows are deleted (saved game not found)
func TestDeleteSavedGameDB_Unhappy_NoRowsAffected(t *testing.T) {
	mock, repo, _ := setupSavedGameRepo(t)

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
// POST/CREATE saved game
/////////////////////////////

// mockSavedGame creates a test SavedGame with encrypted winning player name
func mockSavedGame(encService *encryption.EncryptionService) models.SavedGame {
	encryptedName, _ := encService.Encrypt("Darien")
	hash := []byte("darien_hash")
	
	return models.SavedGame{
		TotalPoints:                1200,
		AveragePoints:              400,
		WinningPlayerNameEncrypted: encryptedName,
		WinningPlayerNameHash:      hash,
		WinningPlayerId:            sql.NullInt32{Int32: 1, Valid: true},
		LocationId:                 1,
		Players: []models.Player{
			{ID: 1, PlayerName: "Darien", Score: 500},
			{ID: 2, PlayerName: "Vicky", Score: 400},
		},
	}
}

// TestAddStandardSavedGame_Happy verifies full successful transaction:
// insert saved game with encrypted name, insert players with encrypted names, commit
func TestAddStandardSavedGame_Happy(t *testing.T) {
	mock, repo, encService := setupSavedGameRepo(t)
	game := mockSavedGame(encService)

	mock.ExpectBegin()
	
	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WithArgs(
			game.TotalPoints,
			game.AveragePoints,
			sqlmock.AnyArg(), // encrypted winning player name
			sqlmock.AnyArg(), // hashed winning player name
			game.WinningPlayerId,
			game.LocationId,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	// Expect encrypted player inserts
	for _, p := range game.Players {
		mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
			WithArgs(
				p.ID,
				sqlmock.AnyArg(), // game ID
				p.Score,
				sqlmock.AnyArg(), // encrypted player name
				sqlmock.AnyArg(), // hashed player name
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	mock.ExpectCommit()

	result := repo.addStandardSavedGame(game)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestAddStandardSavedGame_Unhappy_BeginTxFails verifies error handling when
// transaction begin fails
func TestAddStandardSavedGame_Unhappy_BeginTxFails(t *testing.T) {
	mock, repo, encService := setupSavedGameRepo(t)
	game := mockSavedGame(encService)

	mock.ExpectBegin().WillReturnError(errors.New("tx fail"))

	result := repo.addStandardSavedGame(game)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestAddStandardSavedGame_Unhappy_InsertSavedGameFails verifies rollback when
// inserting the saved game fails before players are inserted
func TestAddStandardSavedGame_Unhappy_InsertSavedGameFails(t *testing.T) {
	mock, repo, encService := setupSavedGameRepo(t)
	game := mockSavedGame(encService)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnError(errors.New("insert fail"))

	mock.ExpectRollback()

	result := repo.addStandardSavedGame(game)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestAddStandardSavedGame_Unhappy_EncryptionFails verifies error handling when
// player name encryption fails during junction table insert
func TestAddStandardSavedGame_Unhappy_EncryptionFails(t *testing.T) {
	mock, repo, encService := setupSavedGameRepo(t)
	game := mockSavedGame(encService)
	
	// Replace encryption service with bad one
	repo.encryptionService = encryption.NewService([]byte("short"))

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	mock.ExpectRollback()

	result := repo.addStandardSavedGame(game)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestAddStandardSavedGame_Unhappy_PlayerInsertFails verifies rollback when
// inserting one of the players into the junction table fails
func TestAddStandardSavedGame_Unhappy_PlayerInsertFails(t *testing.T) {
	mock, repo, encService := setupSavedGameRepo(t)
	game := mockSavedGame(encService)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	// First player succeeds
	mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Second player fails
	mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("player insert fail"))

	mock.ExpectRollback()

	result := repo.addStandardSavedGame(game)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestAddStandardSavedGame_Unhappy_CommitFails verifies error handling when
// transaction commit fails after all inserts succeed
func TestAddStandardSavedGame_Unhappy_CommitFails(t *testing.T) {
	mock, repo, encService := setupSavedGameRepo(t)
	game := mockSavedGame(encService)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	for range game.Players {
		mock.ExpectExec(regexp.QuoteMeta(constants.InsertPlayersForSavedGame)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	mock.ExpectCommit().WillReturnError(errors.New("commit fail"))

	result := repo.addStandardSavedGame(game)

	require.NotNil(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestAddStandardSavedGame_Happy_NoPlayers verifies successful transaction when
// saved game is created with no players in junction table
func TestAddStandardSavedGame_Happy_NoPlayers(t *testing.T) {
	mock, repo, encService := setupSavedGameRepo(t)
	game := mockSavedGame(encService)
	game.Players = []models.Player{} // Empty players

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerSavedGame)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	// No player inserts expected
	mock.ExpectCommit()

	result := repo.addStandardSavedGame(game)

	require.Nil(t, result.Err)
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.NoError(t, mock.ExpectationsWereMet())
}