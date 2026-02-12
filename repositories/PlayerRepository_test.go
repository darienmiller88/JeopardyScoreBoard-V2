package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/encryption"
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

func setupPlayerRepo(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *sqlPlayerRepository, *encryption.EncryptionService) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	
	// Create a test encryption service with a 32-byte key
	testKey := []byte("12345678901234567890123456789012") // 32 bytes for AES-256
	encryptionService := encryption.NewService(testKey)
	
	repo := &sqlPlayerRepository{
		db: sqlxDB,
		encryptionService: encryptionService,
	}

	t.Cleanup(func() {
		db.Close()
	})

	return sqlxDB, mock, repo, encryptionService
}

//////////////////////////////
//POST tests
////////////////////////////

// func TestAddPlayerToLocation_Success(t *testing.T) {
// 	_, mock, repo, encryptionService := setupPlayerRepo(t)

// 	player := models.Player{
// 		PlayerName: "Jane Doe",
// 	}

// 	// Encrypt the name to get what will be stored in DB
// 	encryptedName, err := encryptionService.Encrypt("Jane Doe")
// 	require.NoError(t, err)

// 	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "location_id"}).
// 		AddRow(42, "2024-01-01", "2024-01-01", 1)

// 	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerWithoutTeam)).
// 		WithArgs(encryptedName, sqlmock.AnyArg(), "Elmwood"). // AnyArg for hash since it depends on env var
// 		WillReturnRows(rows)

// 	result := repo.AddPlayerToLocation("Elmwood", player)

// 	require.NoError(t, result.Err)
// 	assert.Equal(t, http.StatusCreated, result.StatusCode)
// 	assert.Equal(t, 42, result.ResultData.ID)
// 	assert.NotEmpty(t, result.ResultData.PlayerNameEncrypted)
// 	assert.NotEmpty(t, result.ResultData.PlayerNameHash)
// }

func TestAddPlayerToLocation_Success(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	player := models.Player{
		PlayerName: "Jane Doe",
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "location_id"}).
		AddRow(42, time.Now(), time.Now(), 1)

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerWithoutTeam)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Elmwood"). // encrypted name, hash, location
		WillReturnRows(rows)

	result := repo.AddPlayerToLocation("Elmwood", player)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, 42, result.ResultData.ID)
	assert.NotEmpty(t, result.ResultData.PlayerNameEncrypted)
	assert.NotEmpty(t, result.ResultData.PlayerNameHash)
}


func TestAddPlayerToLocation_EncryptionError(t *testing.T) {
	_, _, repo, _ := setupPlayerRepo(t)

	// Create a repo with a bad encryption service
	repo.encryptionService = encryption.NewService([]byte("short")) // Invalid key length

	player := models.Player{
		PlayerName: "Jane Doe",
	}

	result := repo.AddPlayerToLocation("Elmwood", player)

	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestAddPlayerToLocation_DatabaseError(t *testing.T) {
	_, mock, repo, encryptionService := setupPlayerRepo(t)

	player := models.Player{
		PlayerName: "Jane Doe",
	}

	encryptedName, _ := encryptionService.Encrypt("Jane Doe")

	mock.ExpectQuery(regexp.QuoteMeta(constants.InsertNewPlayerWithoutTeam)).
		WithArgs(encryptedName, sqlmock.AnyArg(), "Elmwood").
		WillReturnError(fmt.Errorf("database error"))

	result := repo.AddPlayerToLocation("Elmwood", player)

	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

/////////////////////////
// UPDATE tests
////////////////////////

func TestUpdatePlayerName_Happy(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	newName := "Kathya"
	oldName := "Kathy"

	mock.ExpectExec(regexp.QuoteMeta(constants.UpdatePlayerName)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Elmwood"). // All encrypted/hash args
		WillReturnResult(sqlmock.NewResult(0, 1))

	result := repo.UpdatePlayerName(oldName, newName, "Elmwood")

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Nil(t, result.Err)
	assert.Equal(t, newName, result.ResultData.PlayerName)
	assert.NotEmpty(t, result.ResultData.PlayerNameEncrypted)
	assert.NotEmpty(t, result.ResultData.PlayerNameHash)
}

func TestUpdatePlayerName_PlayerNotFound_Unhappy(t *testing.T) {
	_, mock, repo, encryptionService := setupPlayerRepo(t)

	newName := "Kathya"
	oldName := "Nonexistent player"

	encryptedNewName, _ := encryptionService.Encrypt(newName)

	mock.ExpectExec(regexp.QuoteMeta(constants.UpdatePlayerName)).
		WithArgs(encryptedNewName, sqlmock.AnyArg(), sqlmock.AnyArg(), "Elmwood").
		WillReturnResult(sqlmock.NewResult(0, 0))

	result := repo.UpdatePlayerName(oldName, newName, "Elmwood")

	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.NotNil(t, result.Err)
	assert.Contains(t, result.Err.Error(), "could not find player")
}

func TestUpdatePlayerName_EncryptionError(t *testing.T) {
	_, _, repo, _ := setupPlayerRepo(t)

	// Bad encryption service
	repo.encryptionService = encryption.NewService([]byte("short"))

	result := repo.UpdatePlayerName("old", "new", "Elmwood")

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.NotNil(t, result.Err)
}

func TestUpdatePlayerName_ExecError_Unhappy(t *testing.T) {
	_, mock, repo, encryptionService := setupPlayerRepo(t)

	newName := "Kathya"
	oldName := "new person"

	encryptedNewName, _ := encryptionService.Encrypt(newName)

	mock.ExpectExec(regexp.QuoteMeta(constants.UpdatePlayerName)).
		WithArgs(encryptedNewName, sqlmock.AnyArg(), sqlmock.AnyArg(), "Elmwood").
		WillReturnError(fmt.Errorf("database exec error"))

	result := repo.UpdatePlayerName(oldName, newName, "Elmwood")

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.NotNil(t, result.Err)
}




//////////////////////
//DELETE tests
/////////////////////

func TestDeletePlayerName_Success(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	name := "Kathya"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeletePlayer)).
		WithArgs(sqlmock.AnyArg(), "Elmwood"). // hash, location
		WillReturnResult(sqlmock.NewResult(0, 1))

	result := repo.RemovePlayer(name, "Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, name, result.ResultData.PlayerName)
}

func TestDeletePlayerName_OldNameNotFound(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	name := "Kathya"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeletePlayer)).
		WithArgs(sqlmock.AnyArg(), "Elmwood").
		WillReturnResult(sqlmock.NewResult(0, 0))

	result := repo.RemovePlayer(name, "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "could not find player")
}

func TestDeletePlayerName_DatabaseError(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	name := "Kathya"

	mock.ExpectExec(regexp.QuoteMeta(constants.DeletePlayer)).
		WithArgs(sqlmock.AnyArg(), "Elmwood").
		WillReturnError(fmt.Errorf("database error"))

	result := repo.RemovePlayer(name, "Elmwood")

	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}




//////////////////
//GET tests
/////////////////

func TestGetPlayersFromLocation_Ok(t *testing.T) {
	_, mock, repo, encryptionService := setupPlayerRepo(t)

	location := "Elmwood"
	
	// Create encrypted names for test data
	encryptedName1, _ := encryptionService.Encrypt("brent cooper")
	encryptedName2, _ := encryptionService.Encrypt("marky mark")

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"player_name_encrypted",
		"player_name_hash",
		"location_id",
		"team_id",
	}).
		AddRow(1, time.Now(), time.Now(), encryptedName1, []byte("hash1"), 1, nil).
		AddRow(2, time.Now(), time.Now(), encryptedName2, []byte("hash2"), 1, nil)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayersFromLocation)).
		WithArgs(location).
		WillReturnRows(rows)

	result := repo.GetPlayersFromLocation(location)

	require.NoError(t, result.Err)
	assert.Equal(t, 2, len(result.ResultData))
	assert.Equal(t, "brent cooper", result.ResultData[0].PlayerNameDecrypted)
	assert.Equal(t, "marky mark", result.ResultData[1].PlayerNameDecrypted)
}

func TestGetPlayersFromLocation_InvalidLocation(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)
	location := "FakeLocation"

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayersFromLocation)).
		WithArgs(location).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"created_at",
				"updated_at",
				"player_name_encrypted",
				"player_name_hash",
				"location_id",
				"team_id",
			}),
		) //should return no rows

	result := repo.GetPlayersFromLocation(location)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Empty(t, result.ResultData)
}

func TestGetPlayersFromLocation_DatabaseError(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayersFromLocation)).
		WithArgs("Elmwood").
		WillReturnError(fmt.Errorf("database error"))

	result := repo.GetPlayersFromLocation("Elmwood")

	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestGetPlayersFromLocation_DecryptionError(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	// Return invalid encrypted data that can't be decrypted
	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"player_name_encrypted",
		"player_name_hash",
		"location_id",
		"team_id",
	}).
		AddRow(1, "2024-01-01", "2024-01-01", []byte("invalid encrypted data"), []byte("hash1"), 1, nil)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayersFromLocation)).
		WithArgs("Elmwood").
		WillReturnRows(rows)

	result := repo.GetPlayersFromLocation("Elmwood")

	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestGetAllPlayersFromAllLocations_Ok(t *testing.T) {
	_, mock, repo, encryptionService := setupPlayerRepo(t)
	
	encryptedName1, _ := encryptionService.Encrypt("brent cooper")
	encryptedName2, _ := encryptionService.Encrypt("marky mark")
	encryptedName3, _ := encryptionService.Encrypt("dar miller")

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"player_name_encrypted",
		"player_name_hash",
		"location_id",
		"team_id",
	}).
		AddRow(1, time.Now(), time.Now(), encryptedName1, []byte("hash1"), 1, nil).
		AddRow(2, time.Now(), time.Now(), encryptedName2, []byte("hash2"), 1, nil).
		AddRow(3, time.Now(), time.Now(), encryptedName3, []byte("hash3"), 2, nil)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayers)).
		WillReturnRows(rows)

	result := repo.GetAllPlayersFromAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, 3, len(result.ResultData))
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "brent cooper", result.ResultData[0].PlayerNameDecrypted)
	assert.Equal(t, "marky mark", result.ResultData[1].PlayerNameDecrypted)
	assert.Equal(t, "dar miller", result.ResultData[2].PlayerNameDecrypted)
}

func TestGetAllPlayersFromAllLocations_EmptyResult(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"player_name_encrypted",
		"player_name_hash",
		"location_id",
		"team_id",
	})

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayers)).
		WillReturnRows(rows)

	result := repo.GetAllPlayersFromAllLocations()

	require.NoError(t, result.Err)
	assert.Empty(t, result.ResultData)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestGetAllPlayersFromAllLocations_DatabaseError(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayers)).
		WillReturnError(fmt.Errorf("database connection lost"))

	result := repo.GetAllPlayersFromAllLocations()

	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestGetPlayerByName_Success(t *testing.T) {
	_, mock, repo, encryptionService := setupPlayerRepo(t)

	playerName := "John Doe"
	encryptedName, _ := encryptionService.Encrypt(playerName)
	playerHash := []byte("somehash")

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"player_name_encrypted",
		"player_name_hash",
		"location_id",
		"team_id",
	}).
		AddRow(1, time.Now(), time.Now(), encryptedName, playerHash, 1, nil)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetPlayerByName)).
		WithArgs(playerHash).
		WillReturnRows(rows)

	result := repo.GetPlayerByName(playerHash)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, playerName, result.ResultData.PlayerNameDecrypted)
}

func TestGetPlayerByName_NotFound(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	playerHash := []byte("nonexistenthash")

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetPlayerByName)).
		WithArgs(playerHash).
		WillReturnError(sql.ErrNoRows)

	result := repo.GetPlayerByName(playerHash)

	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

func TestGetPlayerByName_DatabaseError(t *testing.T) {
	_, mock, repo, _ := setupPlayerRepo(t)

	playerHash := []byte("somehash")

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetPlayerByName)).
		WithArgs(playerHash).
		WillReturnError(fmt.Errorf("database error"))

	result := repo.GetPlayerByName(playerHash)

	assert.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}