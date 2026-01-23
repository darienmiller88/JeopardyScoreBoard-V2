package repositories

import (
	"JeopardyScoreBoardV2/constants"
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

func setupTeamRepo(t *testing.T) (sqlmock.Sqlmock, *sqlTeamRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := &sqlTeamRepository{db: sqlxDB}

	t.Cleanup(func() {
		db.Close()
	})

	return mock, repo
}

// Happy path: team exists and players are returned
func TestGetTeamWithAllPlayersDB_Happy(t *testing.T) {
	mock, repo := setupTeamRepo(t)
	teamId := 10

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetTeamById)).
		WithArgs(teamId).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "location_id", "created_at", "updated_at"}).
				AddRow(teamId, 3, time.Now(), time.Now()),
		)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayersOnTeam)).
		WithArgs(teamId).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "player_name", "location_id", "team_id", "created_at", "updated_at"}).
				AddRow(1, "Alice", 3, teamId, time.Now(), time.Now()).
				AddRow(2, "Bob", 3, teamId, time.Now(), time.Now()),
		)

	result := repo.GetTeamWithAllPlayersDB(teamId)

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData.Players, 2)
}

// Edge case: team exists but has zero players → still 200 OK
func TestGetTeamWithAllPlayersDB_NoPlayers_Happy(t *testing.T) {
	mock, repo := setupTeamRepo(t)

	teamId := 10

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetTeamById)).
		WithArgs(teamId).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "location_id", "created_at", "updated_at"}).
				AddRow(teamId, 3, time.Now(), time.Now()),
		)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayersOnTeam)).
		WithArgs(teamId).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "player_name", "location_id", "team_id", "created_at", "updated_at"}),
		)

	result := repo.GetTeamWithAllPlayersDB(teamId)

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData.Players, 0)
}

// Happy path: team names returned
func TestGetAllTeamNames_Happy(t *testing.T) {
	mock, repo := setupTeamRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllTeamsByName)).
		WillReturnRows(
			sqlmock.NewRows([]string{"location_name"}).
				AddRow("Elmwood").
				AddRow("Pelham Bay"),
		)

	result := repo.GetAllTeamNames()

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 2)
}



/////////////////////////
//UNHAPPY PATHS
/////////////////////////

// Unhappy: DB error when fetching names → 500
func TestGetAllTeamNames_QueryError_Unhappy(t *testing.T) {
	mock, repo := setupTeamRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllTeamsByName)).
		WillReturnError(errors.New("db error"))

	result := repo.GetAllTeamNames()

	require.Error(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// Unhappy: team id does not exist → 404
func TestGetTeamWithAllPlayersDB_TeamNotFound_Unhappy(t *testing.T) {
	mock, repo := setupTeamRepo(t)

	teamId := 999

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetTeamById)).
		WithArgs(teamId).
		WillReturnError(sql.ErrNoRows)

	result := repo.GetTeamWithAllPlayersDB(teamId)

	require.Error(t, result.Err)
	require.Equal(t, http.StatusNotFound, result.StatusCode)
}

// Unhappy: players query fails → 500
func TestGetTeamWithAllPlayersDB_PlayerQueryError_Unhappy(t *testing.T) {
	mock, repo := setupTeamRepo(t)

	teamId := 10

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetTeamById)).
		WithArgs(teamId).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "location_id", "created_at", "updated_at"}).
				AddRow(teamId, 3, time.Now(), time.Now()),
		)

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetAllPlayersOnTeam)).
		WithArgs(teamId).
		WillReturnError(errors.New("player query failed"))

	result := repo.GetTeamWithAllPlayersDB(teamId)

	require.Error(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// Unhappy: DB error when fetching team → 500
func TestGetTeamWithAllPlayersDB_TeamQueryError_Unhappy(t *testing.T) {
	mock, repo := setupTeamRepo(t)

	teamId := 10

	mock.ExpectQuery(regexp.QuoteMeta(constants.GetTeamById)).
		WithArgs(teamId).
		WillReturnError(errors.New("db failure"))

	result := repo.GetTeamWithAllPlayersDB(teamId)

	require.Error(t, result.Err)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}
