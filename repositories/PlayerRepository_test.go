package repositories

import (
	"JeopardyScoreBoardV2/models"
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRepo(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *sqlPlayerRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := &sqlPlayerRepository{db: sqlxDB}

	t.Cleanup(func() {
		db.Close()
	})

	return sqlxDB, mock, repo
}

//////////////////////////////
//POST tests
////////////////////////////

func TestAddPlayerToLocation_Success(t *testing.T) {
	_, mock, repo := setupRepo(t)

	player := models.Player{
		PlayerName: "Jane Doe",
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(42)

	mock.ExpectQuery(`INSERT INTO players`).
		WithArgs("Jane Doe", "Elmwood").
		WillReturnRows(rows)

	result := repo.AddPlayerToLocation("Elmwood", player)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, 42, result.ResultData.ID)
}

func TestAddPlayerToLocation_FakeLocation(t *testing.T) {
	_, mock, repo := setupRepo(t)

	fakeLocation := "FakeLocation"
	player := models.Player{
		PlayerName: "Jane Doe",
	}

	mock.ExpectQuery(`INSERT INTO players`).
		WithArgs(player.PlayerName, fakeLocation).
		WillReturnError(&pq.Error{
			Code: "23502",
		})

	result := repo.AddPlayerToLocation(fakeLocation, player)

	assert.Equal(t, fmt.Sprintf("no location '%s' found", fakeLocation), result.Err.Error())
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

func TestAddPlayerToLocation_PlayerNameTaken(t *testing.T) {
	_, mock, repo := setupRepo(t)

	location := "Elmwood"
	player := models.Player{
		PlayerName: "Jane Doe",
	}

	mock.ExpectQuery(`INSERT INTO players`).
		WithArgs(player.PlayerName, location).
		WillReturnError(&pq.Error{
			Code: "23505",
		})

	result := repo.AddPlayerToLocation(location, player)

	require.Error(t, result.Err)
	assert.Contains(t, result.Err.Error(), fmt.Sprintf("name %s is already taken", player.PlayerName))
	assert.Equal(t, http.StatusConflict, result.StatusCode)
}

/////////////////////////
// UPDATE tests
////////////////////////

func TestUpdatePlayerName_Ok(t *testing.T) {
	_, mock, repo := setupRepo(t)

	newName := "Kathya"
	oldName := "Kathy"

	mock.ExpectExec(`UPDATE players`).
		WithArgs(newName, oldName).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result := repo.UpdatePlayerName(oldName, newName)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, nil, result.Err)
	assert.Equal(t, newName, result.ResultData.PlayerName)
}
func TestUpdatePlayerName_NewNameTaken(t *testing.T) {
	_, mock, repo := setupRepo(t)

	newName := "Kathya"
	oldName := "Kathy"

	mock.ExpectExec(`UPDATE players`).
		WithArgs(newName, oldName).
		WillReturnError(&pq.Error{
			Code: "23505",
		})

	result := repo.UpdatePlayerName(oldName, newName)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
}

func TestUpdatePlayerName_OldNameNotFound(t *testing.T) {
	_, mock, repo := setupRepo(t)

	newName := "Kathya"
	oldName := "Kathy"

	mock.ExpectExec(`UPDATE players`).
		WithArgs(newName, oldName).
		WillReturnResult(sqlmock.NewResult(0, 0))

	result := repo.UpdatePlayerName(oldName, newName)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

//////////////////////
//DELETE tests
/////////////////////

func TestDeletePlayerName_OldNameNotFound(t *testing.T) {
	_, mock, repo := setupRepo(t)

	name := "Kathya"

	mock.ExpectExec(`DELETE FROM players`).
		WithArgs(name).
		WillReturnResult(sqlmock.NewResult(0, 0))

	result := repo.RemovePlayer(name)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

//////////////////
//GET tests
/////////////////

func TestGetPlayersFromLocation_Ok(t *testing.T) {
	_, mock, repo := setupRepo(t)

	location := "Elmwood"
	rows := sqlmock.NewRows([]string{
		"id",
		"player_name",
	}).AddRow(1, "brent cooper").AddRow(2, "marky mark")

	mock.ExpectQuery(`SELECT players.* FROM players`).
		WithArgs(location).
		WillReturnRows(rows)

	result := repo.GetPlayersFromLocation(location)

	require.NoError(t, result.Err)
	assert.Equal(t, 2, len(result.ResultData))
}

func TestGetPlayersFromLocation_InvalidLocation(t *testing.T) {
	_, mock, repo := setupRepo(t)

	location := "FakeLocation"

	mock.ExpectQuery(`SELECT players.* FROM players`).
		WithArgs(location).
		WillReturnRows(
             sqlmock.NewRows([]string{
                "id",
                "player_name",
                "created_at",
                "updated_at",
                "location_id",
                "team_id",
            }),
        ) //should return no rows

	result := repo.GetPlayersFromLocation(location)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestGetAllPlayersFromAllLocations_Ok(t *testing.T) {
	_, mock, repo := setupRepo(t)

	rows := sqlmock.NewRows([]string{
		"id",
		"player_name",
	}).AddRow(1, "brent cooper").AddRow(2, "marky mark").AddRow(3, "dar miller")

	mock.ExpectQuery(`SELECT \* FROM players`).WillReturnRows(rows)

	result := repo.GetAllPlayersFromAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, 3, len(result.ResultData))
	assert.Equal(t, http.StatusOK, result.StatusCode)
}
