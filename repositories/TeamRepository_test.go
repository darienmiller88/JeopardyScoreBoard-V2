package repositories

import (
	"JeopardyScoreBoardV2/constants"
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
