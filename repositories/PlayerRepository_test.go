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

func TestAddPlayerToLocation_Success(t *testing.T) {
    _, mock, repo := setupRepo(t)

    player := models.Player{
        PlayerName: "Jane Doe",
    }

    rows := sqlmock.NewRows([]string{"id"}).AddRow(42)

    mock.ExpectQuery(`INSERT INTO players`).
        WithArgs("Jane Doe", "New York").
        WillReturnRows(rows)

    result := repo.AddPlayerToLocation("New York", player)

    require.NoError(t, result.Err)
    assert.Equal(t, http.StatusOK, result.StatusCode)
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
