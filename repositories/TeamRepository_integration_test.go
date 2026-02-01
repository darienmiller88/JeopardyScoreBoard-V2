package repositories

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Happy path: retrieves an existing team along with all assigned players
func TestGetTeamWithAllPlayersDB_Integration_Happy(t *testing.T) {
	repo := GetSqlTeamRepository(db)

	// Get an existing team ID
	var teamId int
	err := db.Get(&teamId, `
		SELECT id
		FROM teams
		LIMIT 1
	`)
	require.NoError(t, err)

	result := repo.GetTeamWithAllPlayersDB(teamId)

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Equal(t, teamId, result.ResultData.ID)

	// Players slice should always be non-nil
	require.NotNil(t, result.ResultData.Players)
}

// Happy path: retrieves all team names (location names)
func TestGetAllTeamNames_Integration_Happy(t *testing.T) {
	repo := GetSqlTeamRepository(db)

	result := repo.GetAllTeamNamesDB()

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.NotEmpty(t, result.ResultData)
}

// Happy path: team names exactly match location names
func TestGetAllTeamNames_MatchLocations_Happy(t *testing.T) {
	repo := GetSqlTeamRepository(db)

	var locationNames []string
	err := db.Select(&locationNames, `
		SELECT location_name
		FROM locations
		ORDER BY location_name
	`)
	require.NoError(t, err)

	result := repo.GetAllTeamNamesDB()

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.ElementsMatch(t, locationNames, result.ResultData)
}

// Edge case: no teams exist → returns empty slice with 200 OK
func TestGetAllTeamNames_NoTeams_Happy(t *testing.T) {
	repo := GetSqlTeamRepository(db)

	// Temporarily clear teams
	_, err := db.Exec("DELETE FROM teams")
	require.NoError(t, err)

	t.Cleanup(func() {
		// Re-seed teams from locations
		_, _ = db.Exec(`
			INSERT INTO teams (location_id)
			SELECT id FROM locations
			ON CONFLICT DO NOTHING
		`)
	})

	result := repo.GetAllTeamNamesDB()

	require.NoError(t, result.Err)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Len(t, result.ResultData, 0)
}

// Unhappy path: team ID does not exist → 404
func TestGetTeamWithAllPlayersDB_Integration_TeamNotFound_Unhappy(t *testing.T) {
	repo := GetSqlTeamRepository(db)

	result := repo.GetTeamWithAllPlayersDB(999999)

	require.Error(t, result.Err)
	require.Equal(t, http.StatusNotFound, result.StatusCode)
}
