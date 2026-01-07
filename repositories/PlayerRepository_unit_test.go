package repositories

import (
	"JeopardyScoreBoardV2/models"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
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