package services

import (
	"fmt"
	"net/http"
	"testing"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlayerRepository struct {
	playerResult  models.Result[models.Player]
	playersResult models.Result[[]models.Player]
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

type mockLocationRepositoryForPlayer struct {
	locationResult models.Result[string]
}

func (m *mockLocationRepositoryForPlayer) GetLocation(locationName string) models.Result[string] {
	return m.locationResult
}

func (m *mockLocationRepositoryForPlayer) GetLocationById(locationId int)  models.Result[models.Location]{
	return models.Result[models.Location]{}
}

func (m *mockLocationRepositoryForPlayer) GetAllLocations() models.Result[[]string] {
	return models.Result[[]string]{}
}


////////////////////
// CREATE/POST tests
////////////////////

// TestAddPlayer_Ok verifies successful player addition when all validations pass
func TestAddPlayer_Ok(t *testing.T) {
	firstName := "Jane"
	lastName := "Doe"
	mockPlayerRepo := &mockPlayerRepository{
		getPlayerByNameFunc: func(playerName string) models.Result[models.Player] {
			// First call from isPlayerNameTaken - name not found (good)
			return utils.GetResult(
				fmt.Errorf("player not found"),
				http.StatusNotFound,
				models.Player{},
			)
		},
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{
				FirstName: firstName,
				LastName: lastName,
			},
			StatusCode: http.StatusCreated, // ← Change this to 201
		},
	}
	
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.AddPlayerToLocation("Elmwood", firstName, lastName)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, "Jane Doe", result.ResultData.PlayerName)
}

// TestAddPlayer_NameTooShort verifies validation fails for names under 4 characters
func TestAddPlayer_FirstNameTooShort(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	mockLocationRepo := &mockLocationRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}

	result := service.AddPlayerToLocation("Elmwood", "J", "liberman")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestAddPlayer_NameTooLong verifies validation fails for first names over 20 characters
func TestAddPlayer_FirstNameTooLong(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	mockLocationRepo := &mockLocationRepository{}
	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}

	result := service.AddPlayerToLocation("Elmwood", "Joedcsxevrgvfsxergtdwxertgfwsxdgtrvwsxdertgvcevrtbgfcdwextra", "regular")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestAddPlayer_NameMustHaveTwoParts verifies validation requires exactly two name parts
func TestAddPlayer_NameMustHaveTwoParts(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	mockLocationRepo := &mockLocationRepositoryForPlayer{}
	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}

	// Single name
	result := service.AddPlayerToLocation("Elmwood", "Cheryl")
	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "two parts")

	// More than two names
	result = service.AddPlayerToLocation("Elmwood", "This name has more than two parts")
	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "two parts")
}

// TestAddPlayer_LocationNotFound verifies error when location doesn't exist
func TestAddPlayer_LocationNotFound(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			Err:        fmt.Errorf("location not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.AddPlayerToLocation("NonExistent", "Jane Doe")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestAddPlayer_NameAlreadyTaken verifies error when player name already exists
func TestAddPlayer_NameAlreadyTaken(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{PlayerName: "Jane Doe"},
			StatusCode: http.StatusOK, // Player exists
		},
	}
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.AddPlayerToLocation("Elmwood", "Jane Doe")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "is taken")
}

// TestAddPlayer_DatabaseError verifies error handling for database failures
func TestAddPlayer_DatabaseError(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("database connection lost"),
			StatusCode: http.StatusInternalServerError,
		},
	}
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.AddPlayerToLocation("Elmwood", "Jane Doe")

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
	count := 0

	mockPlayerRepo := &mockPlayerRepository{
		getPlayerByNameFunc: func(playerName string) models.Result[models.Player] {
			count++

			// First call from doesPlayerExist - name found (good)
			if count == 1 {
				return utils.GetResult(nil, http.StatusOK, models.Player{})
			}

			//Second call from isPlayerTaken - new name not found (good)
			return utils.GetResult(nil, int(http.StatusNotFound), models.Player{})
		},
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{PlayerName: "Bob Melendez"},
			StatusCode: http.StatusOK, 
		},
	}
	
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.UpdatePlayerName("Alice Twilight", "Bob Melendez", "Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "Bob Melendez", result.ResultData.PlayerName)
}

// TestUpdatePlayerName_Service_InvalidNewName verifies validation fails for invalid new name
func TestUpdatePlayerName_Service_InvalidNewName(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	mockLocationRepo := &mockLocationRepositoryForPlayer{}
	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}

	result := service.UpdatePlayerName("Alice Twilight", "Bob", "Elmwood") // too short

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
}

// TestUpdatePlayerName_Service_NameMustHaveTwoParts verifies new name must have two parts
func TestUpdatePlayerName_Service_NameMustHaveTwoParts(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	mockLocationRepo := &mockLocationRepositoryForPlayer{}
	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}

	result := service.UpdatePlayerName("Alice Twilight", "Margaret", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "two parts")
}

// TestUpdatePlayerName_Service_SameName verifies error when old and new names are identical
func TestUpdatePlayerName_Service_SameName(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	mockLocationRepo := &mockLocationRepositoryForPlayer{}
	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}

	result := service.UpdatePlayerName("Jane Doe", "Jane Doe", "Elmwood")

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
	mockLocationRepo := &mockLocationRepository{}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.UpdatePlayerName("Ghost Player", "Jane Doe", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestUpdatePlayerName_Service_NewNameAlreadyTaken verifies error when new name is taken
func TestUpdatePlayerName_Service_NewNameAlreadyTaken(t *testing.T) {
	callCount := 0

	//UpdatePlayerName calls GetPlayerByName(playerName string) twice: once to verify if the old name exists,
	//and another to verify if the new name is taken. If the second call returns a 200, the name is taken
	//and the update cannot happen.
	mockPlayerRepo := &mockPlayerRepository{
		getPlayerByNameFunc: func(playerName string) models.Result[models.Player] {
			callCount++

			if callCount == 1 {
				// First call: old name exists (OK)
				return utils.GetResult(nil, http.StatusOK, models.Player{PlayerName: "Alice Twilight"})
			}

			// Second call: new name already taken (Conflict)
			return utils.GetResult(nil, http.StatusOK, models.Player{PlayerName: "Bob Melendez"})
		},
	}

	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}

	result := service.UpdatePlayerName("Alice Twilight", "Bob Melendez", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "is taken")
}

// TestUpdatePlayerName_Service_LocationNotFound verifies error when location doesn't exist
func TestUpdatePlayerName_Service_LocationNotFound(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{PlayerName: "Alice Twilight"},

			//to update the name, a 404 must be thrown to signal is isn't there.
			StatusCode: http.StatusNotFound,
		},
	}
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			Err:        fmt.Errorf("location not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.UpdatePlayerName("Alice Twilight", "Bob Melendez", "NonExistent")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestUpdatePlayerName_Service_RepoError verifies error handling for database failures
func TestUpdatePlayerName_Service_RepoError(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("db down"),
			StatusCode: http.StatusInternalServerError,
		},
	}
	mockLocationRepo := &mockLocationRepositoryForPlayer{}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.UpdatePlayerName("Alice Twilight", "Bob Melendez", "Elmwood")

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
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
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
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.RemovePlayer("Ghost Casper", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestRemovePlayer_Service_LocationNotFound verifies error when location doesn't exist
func TestRemovePlayer_Service_LocationNotFound(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{}
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			Err:        fmt.Errorf("location not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
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
	mockLocationRepo := &mockLocationRepositoryForPlayer{
		locationResult: models.Result[string]{
			ResultData: "Elmwood",
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{
		PlayerRepository:   mockPlayerRepo,
		LocationRepository: mockLocationRepo,
	}
	result := service.RemovePlayer("Jane Doe", "Elmwood")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

////////////////////
// Helper Method Tests
////////////////////

// TestDoesPlayerExist_PlayerExists verifies helper returns success when player exists
func TestDoesPlayerExist_PlayerExists(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{PlayerName: "Jane Doe"},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.doesPlayerExist("Jane Doe")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

// TestDoesPlayerExist_PlayerNotFound verifies helper returns 404 when player doesn't exist
func TestDoesPlayerExist_PlayerNotFound(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("player not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.doesPlayerExist("Ghost Player")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestIsPlayerNameTaken_NameIsTaken verifies helper returns conflict when name exists
func TestIsPlayerNameTaken_NameIsTaken(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			ResultData: models.Player{PlayerName: "Jane Doe"},
			StatusCode: http.StatusOK,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.isPlayerNameTaken("Jane Doe")

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "is taken")
}

// TestIsPlayerNameTaken_NameIsAvailable verifies helper returns success when name is available
func TestIsPlayerNameTaken_NameIsAvailable(t *testing.T) {
	mockPlayerRepo := &mockPlayerRepository{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("player not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &PlayerServiceImpl{PlayerRepository: mockPlayerRepo}
	result := service.isPlayerNameTaken("Available Name")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}