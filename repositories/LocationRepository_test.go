package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLocationRepo(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *sqlLocationRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := &sqlLocationRepository{db: sqlxDB}

	t.Cleanup(func() {
		db.Close()
	})

	return sqlxDB, mock, repo
}

func TestGetAllLocations_Ok(t *testing.T) {
	_, mock, repo := setupLocationRepo(t)

	rows := sqlmock.NewRows([]string{"location_name"}).
		AddRow("Elmwood").
		AddRow("Pelham Bay").
		AddRow("Flushing")

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllLocations)).WillReturnRows(rows)

	result := repo.GetAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 3)
	assert.Equal(t, []string{"Elmwood", "Pelham Bay", "Flushing"}, result.ResultData)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllLocations_Empty(t *testing.T) {
	_, mock, repo := setupLocationRepo(t)

	rows := sqlmock.NewRows([]string{"location_name"})

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllLocations)).WillReturnRows(rows)

	result := repo.GetAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Empty(t, result.ResultData)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllLocations_DbError(t *testing.T) {
	_, mock, repo := setupLocationRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllLocations)).
		WillReturnError(fmt.Errorf("db blew up"))

	result := repo.GetAllLocations()

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Empty(t, result.ResultData)

	require.NoError(t, mock.ExpectationsWereMet())
}
