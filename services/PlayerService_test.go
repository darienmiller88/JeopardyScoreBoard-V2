package services

import (
	"fmt"
	"net/http"
	"testing"

	"JeopardyScoreBoardV2/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlayerRepository struct {
	playerResult        models.Result[models.Player]
	playersResult       models.Result[[]models.Player]
	getPlayerByNameFunc func(string) models.Result[models.Player] // Add this function field
}

func (m *mockPlayerRepository) UpdatePlayerName(oldPlayerName string, newPlayerName string, locationName string) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerRepository) AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerRepository) RemovePlayer(playerName string, locationName string) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerRepository) GetPlayersFromLocation(locationName string) models.Result[[]models.Player] {
	return m.playersResult
}

func (m *mockPlayerRepository) GetAllPlayersFromAllLocations() models.Result[[]models.Player] {
	return m.playersResult
}

func (m *mockPlayerRepository) GetPlayersByNames(players []string) models.Result[[]models.Player] {
	return m.playersResult
}

func (m *mockPlayerRepository) GetPlayerByName(playerName string) models.Result[models.Player] {
	// Use the function field if provided, otherwise use default behavior
	if m.getPlayerByNameFunc != nil {
		return m.getPlayerByNameFunc(playerName)
	}

	return m.playerResult
}

////////////////////
// CREATE/POST tests
////////////////////

// TestAddPlayer_Ok verifies successful player addition when all validations pass
func TestAddPlayer_Ok(t *testing.T) {
	firstName := "Jane"
	lastName := "Doe"
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{
				PlayerName: "Jane Doe",
			},
			StatusCode: http.StatusCreated,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}
	result := service.AddPlayerToLocation("Elmwood", firstName, lastName)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, "Jane Doe", result.ResultData.PlayerName)
}

// TestAddPlayer_NameTooShort verifies validation fails for first names under 2 characters
func TestAddPlayer_FirstNameTooShort(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.AddPlayerToLocation("Elmwood", "J", "liberman")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestAddPlayer_NameTooLong verifies validation fails for first names over 20 characters
func TestAddPlayer_FirstNameTooLong(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.AddPlayerToLocation("Elmwood", "Joedcsxevrgvfsxergtdwxertgfwsxdgtrvwsxdertgvcevrtbgfcdwextra", "regular")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestAddPlayer_LastNameTooShort verifies validation fails for last names too short
func TestAddPlayer_LastNameTooShort(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.AddPlayerToLocation("Elmwood", "liberman", "K")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestAddPlayer_LastNameTooLong verifies validation fails for last names over 20 characters
func TestAddPlayer_LastNameTooLong(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.AddPlayerToLocation("Elmwood", "march", "Joedcsxevrgvfsxergtdwxertgfwsxdgtrvwsxdertgvcevrtbgfcdwextra")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestAddPlayer_FirstNameMustHaveOnePart verifies validation requires exactly one part
func TestAddPlayer_FirstNameMustHaveOnePart(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	// Single name
	result := service.AddPlayerToLocation("Elmwood", "marky maek", "kltveu")
	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "1 part")
}

// TestAddPlayer_FirstNameMustHaveOnePart verifies validation requires exactly one part
func TestAddPlayer_LastNameMustHaveOnePart(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	// Single name
	result := service.AddPlayerToLocation("Elmwood", "kat", "marky mark")
	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "1 part")
}

// TestAddPlayer_LocationNotFound verifies error when location doesn't exist
func TestAddPlayer_LocationNotFound(t *testing.T) {
	service := &PlayerServiceImpl{
		PlayerRepository: &mockPlayerRepository{
			playerResult: models.Result[models.Player]{
				Err: fmt.Errorf("location not found"),
				StatusCode: http.StatusNotFound,
			},
		},
	}
	result := service.AddPlayerToLocation("NonExistent", "Jane", "Doe")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestAddPlayer_NameAlreadyTaken verifies error when player name already exists
func TestAddPlayer_NameAlreadyTaken(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err: fmt.Errorf("player taken"),
			StatusCode: http.StatusConflict, // Player exists
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}
	result := service.AddPlayerToLocation("Elmwood", "Jane", "Doe")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "taken")
}

// TestAddPlayer_DatabaseError verifies error handling for database failures
func TestAddPlayer_DatabaseError(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("database connection lost"),
			StatusCode: http.StatusInternalServerError,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}
	result := service.AddPlayerToLocation("Elmwood", "Jane", "Doe")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

////////////////////
// READ/GET tests
////////////////////

// TestGetAllPlayersFromAllLocations_Ok verifies successful retrieval of all players
func TestGetAllPlayersFromAllLocations_Ok(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playersResult: models.Result[[]models.Player]{
			ResultData: []models.Player{
				{PlayerName: "Jane Doe"},
				{PlayerName: "John Smith"},
			},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.GetAllPlayersFromAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 2)
}

// TestGetAllPlayersFromAllLocations_EmptyResult verifies handling of no players
func TestGetAllPlayersFromAllLocations_EmptyResult(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playersResult: models.Result[[]models.Player]{
			ResultData: []models.Player{},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.GetAllPlayersFromAllLocations()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Empty(t, result.ResultData)
}

// TestGetAllPlayersFromAllLocations_DatabaseError verifies error handling
func TestGetAllPlayersFromAllLocations_DatabaseError(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playersResult: models.Result[[]models.Player]{
			Err:        fmt.Errorf("database error"),
			StatusCode: http.StatusInternalServerError,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.GetAllPlayersFromAllLocations()

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestGetPlayersFromLocation_Ok verifies successful retrieval of players from location
func TestGetPlayersFromLocation_Ok(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playersResult: models.Result[[]models.Player]{
			ResultData: []models.Player{
				{PlayerName: "Jane Doe"},
			},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.GetPlayersFromLocation("Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 1)
}

// TestGetPlayersFromLocation_EmptyResult verifies handling when location has no players
func TestGetPlayersFromLocation_EmptyResult(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playersResult: models.Result[[]models.Player]{
			ResultData: []models.Player{},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.GetPlayersFromLocation("Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Empty(t, result.ResultData)
}

// TestGetPlayersFromLocation_RepoError verifies error handling for database failures
func TestGetPlayersFromLocation_RepoError(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playersResult: models.Result[[]models.Player]{
			Err:        fmt.Errorf("db exploded"),
			StatusCode: http.StatusInternalServerError,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.GetPlayersFromLocation("Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

////////////////////
// UPDATE/PUT tests
////////////////////

// TestUpdatePlayerName_Service_Ok verifies successful player name update
func TestUpdatePlayerName_Service_Ok(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{PlayerName: "Bob Melendez"},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}
	result := service.UpdatePlayerName("Alice Twilight", "Bob", "Melendez", "Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "Bob Melendez", result.ResultData.PlayerName)
}

// TestUpdatePlayerName_Service_FirstNameTooShort verifies validation fails for short new first name
func TestUpdatePlayerName_Service_FirstNameTooShort(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.UpdatePlayerName("Alice Twilight", "m", "Tyler", "Elmwood") // too short

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestUpdatePlayerName_Service_LastNameTooShort verifies validation fails for short new last name
func TestUpdatePlayerName_Service_LastNameTooShort(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.UpdatePlayerName("Alice Twilight", "Bob", "k", "Elmwood") // too short

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestUpdatePlayerName_Service_FirstNameTooLong verifies validation fails for long new first name
func TestUpdatePlayerName_Service_FirstNameTooLong(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.UpdatePlayerName("Alice Twilight", "gfthyghugyuyftghjgftyghjugyftyhub", "Tyler", "Elmwood") // too short

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestUpdatePlayerName_Service_LastNameTooLong verifies validation fails for long new last name
func TestUpdatePlayerName_Service_LastNameTooLong(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.UpdatePlayerName("Alice Twilight", "Bob", "Tyljkljkjhkhuitybunnnner", "Elmwood") // too short

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestUpdatePlayerName_Service_FirstNameMustHaveOnePart verifies new first name must have one part
func TestUpdatePlayerName_Service_FirstNameMustHaveOnePart(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.UpdatePlayerName("Alice Twilight", "Margaret 2ho", "blahblah", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "1 part")
}

// TestUpdatePlayerName_Service_LastNameMustHaveOnePart verifies new last name must have one part
func TestUpdatePlayerName_Service_LastNameMustHaveOnePart(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.UpdatePlayerName("Alice Twilight", "Margaret", "  blahblah lbas", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "1 part")
}

// TestUpdatePlayerName_Service_SameName verifies error when old and new names are identical
func TestUpdatePlayerName_Service_SameName(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.UpdatePlayerName("Jane Doe", "Jane", "Doe", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "must be different")
}

// TestUpdatePlayerName_Service_OldNameDoesNotExist verifies error when old name doesn't exist
func TestUpdatePlayerName_Service_OldNameDoesNotExist(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("player not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}
	result := service.UpdatePlayerName("Ghost Player", "Jane", "Doe", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestUpdatePlayerName_Service_NewNameAlreadyTaken verifies error when new name is taken
func TestUpdatePlayerName_Service_NewNameAlreadyTaken(t *testing.T) {
	service := &PlayerServiceImpl{
		PlayerRepository: &mockPlayerRepository{
			playerResult: models.Result[models.Player]{
				StatusCode: http.StatusConflict,
				Err:        fmt.Errorf("name taken"),
			},
		},
	}

	result := service.UpdatePlayerName("Alice Twilight", "Bob", "Melendez", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "taken")
}

// TestUpdatePlayerName_Service_LocationNotFound verifies error when location doesn't exist
func TestUpdatePlayerName_Service_LocationNotFound(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err: fmt.Errorf("location not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}
	result := service.UpdatePlayerName("Alice Twilight", "Bob", "Melendez", "NonExistent")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "location not found")
}

// TestUpdatePlayerName_Service_RepoError verifies error handling for database failures
func TestUpdatePlayerName_Service_RepoError(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("db down"),
			StatusCode: http.StatusInternalServerError,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}
	result := service.UpdatePlayerName("Alice Twilight", "Bob", "Melendez", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

////////////////////
// DESTROY/DELETE tests
////////////////////

// TestRemovePlayer_Service_Ok verifies successful player deletion
func TestRemovePlayer_Service_Ok(t *testing.T) {
	playerDeleted := "Jane Doe"
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{PlayerName: playerDeleted},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.RemovePlayer(playerDeleted, "Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, playerDeleted, result.ResultData.PlayerName)
}

// TestRemovePlayer_Service_PlayerNotFound verifies error when player doesn't exist
func TestRemovePlayer_Service_PlayerNotFound(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("player not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}
	result := service.RemovePlayer("Ghost Casper", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestRemovePlayer_Service_LocationNotFound verifies error when location doesn't exist
func TestRemovePlayer_Service_LocationNotFound(t *testing.T) {
	service := &PlayerServiceImpl{
		PlayerRepository: &mockPlayerRepository{
			playerResult: models.Result[models.Player]{
				Err: fmt.Errorf("location not found"),
				StatusCode: http.StatusNotFound,
			},
		},
	}
	result := service.RemovePlayer("Jane Doe", "NonExistent")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestRemovePlayer_Service_DatabaseError verifies error handling for database failures
func TestRemovePlayer_Service_DatabaseError(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("database connection lost"),
			StatusCode: http.StatusInternalServerError,
		},
	}
	service := &PlayerServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}
	result := service.RemovePlayer("Jane Doe", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}
