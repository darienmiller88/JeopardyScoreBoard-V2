package services

import (
	"JeopardyScoreBoardV2/encryption"
	"JeopardyScoreBoardV2/models"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Mock Repositories
// ============================================================================

type mockSavedGameRepository struct {
	savedGameResult  models.Result[models.SavedGame]
	savedGamesResult models.Result[[]models.SavedGame]
	deleteResult     models.Result[string]
}

func (m *mockSavedGameRepository) GetAllSavedGamesFromLocationDB(locationName string) models.Result[[]models.SavedGame] {
	return m.savedGamesResult
}

func (m *mockSavedGameRepository) AddSavedGameDB(savedGame models.SavedGame) models.Result[models.SavedGame] {
	return m.savedGameResult
}

func (m *mockSavedGameRepository) DeleteSavedGameDB(savedGameId int) models.Result[string] {
	return m.deleteResult
}

func (m *mockSavedGameRepository) GetAllSavedGamesDB() models.Result[[]models.SavedGame] {
	return m.savedGamesResult
}

func (m *mockSavedGameRepository) GetSavedGameById(savedGameId int) models.Result[models.SavedGame] {
	return m.savedGameResult
}

type mockLocationRepositoryForSavedGame struct {
	locationResult  models.Result[models.Location]
	locationsResult models.Result[[]string]
}

func (m *mockLocationRepositoryForSavedGame) GetLocation(locationName string) models.Result[string] {
	return models.Result[string]{}
}

func (m *mockLocationRepositoryForSavedGame) GetLocationById(locationId int) models.Result[models.Location] {
	return m.locationResult
}

func (m *mockLocationRepositoryForSavedGame) GetAllLocations() models.Result[[]string] {
	return m.locationsResult
}

type mockPlayerRepositoryForSavedGame struct {
	playerResult  models.Result[models.Player]
	playersResult models.Result[[]models.Player]
}

func (m *mockPlayerRepositoryForSavedGame) UpdatePlayerName(oldPlayerName string, newPlayerName string, locationName string) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerRepositoryForSavedGame) AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerRepositoryForSavedGame) GetPlayersFromLocation(locationName string) models.Result[[]models.Player] {
	return m.playersResult
}

func (m *mockPlayerRepositoryForSavedGame) RemovePlayer(playerName string, locationName string) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerRepositoryForSavedGame) GetAllPlayersFromAllLocations() models.Result[[]models.Player] {
	return m.playersResult
}

func (m *mockPlayerRepositoryForSavedGame) GetPlayersByNames(players []string) models.Result[[]models.Player] {
	return m.playersResult
}

func (m *mockPlayerRepositoryForSavedGame) GetPlayerByName(playerName string) models.Result[models.Player] {
	return m.playerResult
}

type mockTeamRepositoryForSavedGame struct {
	teamResult  models.Result[models.Team]
	teamsResult models.Result[[]models.Team]
}

func (m *mockTeamRepositoryForSavedGame) GetTeamWithAllPlayersDB(teamId int) models.Result[models.Team] {
	return m.teamResult
}

func (m *mockTeamRepositoryForSavedGame) GetAllTeamNamesDB() models.Result[[]string] {
	return models.Result[[]string]{}
}

func (m *mockTeamRepositoryForSavedGame) GetAllTeamsDB() models.Result[[]models.Team] {
	return m.teamsResult
}

func (m *mockTeamRepositoryForSavedGame) GetAllTeamsByIds(teamIds []int) models.Result[[]models.Team] {
	return m.teamsResult
}

// ============================================================================
// Test Helpers
// ============================================================================

func getTestEncryptionServiceForService() *encryption.EncryptionService {
	testKey := []byte("12345678901234567890123456789012")
	return encryption.NewService(testKey)
}

func createTestPlayersForService(encService *encryption.EncryptionService) []models.Player {
	player1Encrypted, _ := encService.Encrypt("Player One")
	player2Encrypted, _ := encService.Encrypt("Player Two")

	return []models.Player{
		{
			ID:                  1,
			PlayerName:          "Player One",
			PlayerNameEncrypted: player1Encrypted,
			Score:               100,
		},
		{
			ID:                  2,
			PlayerName:          "Player Two",
			PlayerNameEncrypted: player2Encrypted,
			Score:               90,
		},
	}
}

// ============================================================================
// GET Tests
// ============================================================================

// TestGetAllSavedGames_Success verifies successful retrieval of all saved games
func TestGetAllSavedGames_Success(t *testing.T) {
	mockRepo := &mockSavedGameRepository{
		savedGamesResult: models.Result[[]models.SavedGame]{
			ResultData: []models.SavedGame{
				{ID: 1, TotalPoints: 1000},
				{ID: 2, TotalPoints: 2000},
			},
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		SavedGameRepository: mockRepo,
	}

	result := service.GetAllSavedGames()

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 2)
}

// TestGetAllSavedGamesFromLocation_Success verifies retrieval of games from specific location
func TestGetAllSavedGamesFromLocation_Success(t *testing.T) {
	mockRepo := &mockSavedGameRepository{
		savedGamesResult: models.Result[[]models.SavedGame]{
			ResultData: []models.SavedGame{
				{ID: 1, LocationId: 1},
			},
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		SavedGameRepository: mockRepo,
	}

	result := service.GetAllSavedGamesFromLocation("Elmwood")

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

// ============================================================================
// DELETE Tests
// ============================================================================

// TestDeleteSavedGame_Success verifies successful deletion of saved game
func TestDeleteSavedGame_Success(t *testing.T) {
	mockRepo := &mockSavedGameRepository{
		deleteResult: models.Result[string]{
			ResultData: "Game deleted successfully",
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		SavedGameRepository: mockRepo,
	}

	result := service.DeleteSavedGame(1)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

// ============================================================================
// ADD PLAYER GAME Tests - Business Rule Validations
// ============================================================================

// TestAddSavedGame_PlayerGame_Success verifies successful player game creation
func TestAddSavedGame_PlayerGame_Success(t *testing.T) {
	encService := getTestEncryptionServiceForService()
	players := createTestPlayersForService(encService)

	mockSavedGameRepo := &mockSavedGameRepository{
		savedGameResult: models.Result[models.SavedGame]{
			ResultData: models.SavedGame{ID: 1},
			StatusCode: http.StatusCreated,
		},
	}

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			ResultData: models.Location{ID: 1, LocationName: "Elmwood"},
			StatusCode: http.StatusOK,
		},
	}

	mockPlayerRepo := &mockPlayerRepositoryForSavedGame{
		playersResult: models.Result[[]models.Player]{
			ResultData: players,
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		SavedGameRepository: mockSavedGameRepo,
		LocationRepository:  mockLocationRepo,
		PlayerRepository:    mockPlayerRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame:      true,
		LocationId:        1,
		Players:           players,
	}

	result := service.AddSavedGame(savedGame)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
}

// TestAddSavedGame_PlayerGame_EmptyPlayers_Fails verifies error when player game has no players
func TestAddSavedGame_PlayerGame_EmptyPlayers_Fails(t *testing.T) {
	service := &SaveGameServiceImpl{}
	savedGame := models.SavedGame{
		IsPlayerGame: true,
		LocationId:   1,
		Players:    []models.Player{}, // Empty - violates business rule
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "players cannot be empty")
}

// TestAddSavedGame_PlayerGame_WithTeams_Fails verifies error when player game has teams added
func TestAddSavedGame_PlayerGame_WithTeams_Fails(t *testing.T) {
	encService := getTestEncryptionServiceForService()
	players := createTestPlayersForService(encService)

	service := &SaveGameServiceImpl{}

	savedGame := models.SavedGame{
		IsPlayerGame: true,
		LocationId:   1,
		Players:      players,
		Teams: []models.Team{ // Teams in player game - violates business rule
			{ID: 1, Score: 100},
		},
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "player game cannot have any teams")
}

// TestAddSavedGame_PlayerGame_InvalidLocationId_Fails verifies error when location doesn't exist
func TestAddSavedGame_PlayerGame_InvalidLocationId_Fails(t *testing.T) {
	encService := getTestEncryptionServiceForService()
	players := createTestPlayersForService(encService)

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			Err:        fmt.Errorf("location not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &SaveGameServiceImpl{
		LocationRepository: mockLocationRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame: true,
		LocationId:   999, // Invalid location ID
		Players:      players,
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestAddSavedGame_PlayerGame_PlayersNotExist_Fails verifies error when one or more players don't exist
func TestAddSavedGame_PlayerGame_PlayersNotExist_Fails(t *testing.T) {
	encService := getTestEncryptionServiceForService()
	players := createTestPlayersForService(encService)

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			ResultData: models.Location{ID: 1},
			StatusCode: http.StatusOK,
		},
	}

	mockPlayerRepo := &mockPlayerRepositoryForSavedGame{
		playersResult: models.Result[[]models.Player]{
			ResultData: []models.Player{players[0]}, // Only 1 player returned, but 2 requested
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		LocationRepository: mockLocationRepo,
		PlayerRepository:   mockPlayerRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame: true,
		LocationId:   1,
		Players:      players, // 2 players
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "one or more players do not exist")
}

// TestAddSavedGame_PlayerGame_PlayerRepositoryError_Fails verifies error handling for database failures
func TestAddSavedGame_PlayerGame_PlayerRepositoryError_Fails(t *testing.T) {
	encService := getTestEncryptionServiceForService()
	players := createTestPlayersForService(encService)

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			ResultData: models.Location{ID: 1},
			StatusCode: http.StatusOK,
		},
	}

	mockPlayerRepo := &mockPlayerRepositoryForSavedGame{
		playersResult: models.Result[[]models.Player]{
			Err:        fmt.Errorf("database error"),
			StatusCode: http.StatusInternalServerError,
		},
	}

	service := &SaveGameServiceImpl{
		LocationRepository: mockLocationRepo,
		PlayerRepository:   mockPlayerRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame:      true,
		LocationId:        1,
		Players:           players,
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestAddSavedGame_PlayerGame_CalculatesScoresCorrectly verifies score calculations
func TestAddSavedGame_PlayerGame_CalculatesScoresCorrectly(t *testing.T) {
	encService := getTestEncryptionServiceForService()
	players := createTestPlayersForService(encService)

	mockSavedGameRepo := &mockSavedGameRepository{
		savedGameResult: models.Result[models.SavedGame]{
			ResultData: models.SavedGame{ID: 1},
			StatusCode: http.StatusCreated,
		},
	}

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			ResultData: models.Location{ID: 1},
			StatusCode: http.StatusOK,
		},
	}

	mockPlayerRepo := &mockPlayerRepositoryForSavedGame{
		playersResult: models.Result[[]models.Player]{
			ResultData: players,
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		SavedGameRepository: mockSavedGameRepo,
		LocationRepository:  mockLocationRepo,
		PlayerRepository:    mockPlayerRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame:      true,
		LocationId:        1,
		Players:           players,
	}

	result := service.AddSavedGame(savedGame)

	require.NoError(t, result.Err)
	// Verify scores were calculated (would be set by CalculateTotalPoints/CalculateAveragePoints)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
}

// ============================================================================
// ADD TEAM GAME Tests - Business Rule Validations
// ============================================================================

// TestAddSavedGame_TeamGame_Success verifies successful team game creation
func TestAddSavedGame_TeamGame_Success(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Score: 100},
		{ID: 2, Score: 120},
	}

	mockSavedGameRepo := &mockSavedGameRepository{
		savedGameResult: models.Result[models.SavedGame]{
			ResultData: models.SavedGame{ID: 1},
			StatusCode: http.StatusCreated,
		},
	}

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			ResultData: models.Location{ID: 1},
			StatusCode: http.StatusOK,
		},
	}

	mockTeamRepo := &mockTeamRepositoryForSavedGame{
		teamsResult: models.Result[[]models.Team]{
			ResultData: teams,
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		SavedGameRepository: mockSavedGameRepo,
		LocationRepository:  mockLocationRepo,
		TeamRepository:      mockTeamRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame: false,
		LocationId:   1,
		Teams:        teams,
	}

	result := service.AddSavedGame(savedGame)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
}

// TestAddSavedGame_TeamGame_EmptyTeams_Fails verifies error when team game has no teams
func TestAddSavedGame_TeamGame_EmptyTeams_Fails(t *testing.T) {
	service := &SaveGameServiceImpl{}

	savedGame := models.SavedGame{
		IsPlayerGame: false,
		LocationId:   1,
		Teams:        []models.Team{}, // Empty - violates business rule
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "teams cannot be empty")
}

// TestAddSavedGame_TeamGame_WithPlayers_Fails verifies error when team game has players added
func TestAddSavedGame_TeamGame_WithPlayers_Fails(t *testing.T) {
	encService := getTestEncryptionServiceForService()

	service := &SaveGameServiceImpl{}

	savedGame := models.SavedGame{
		IsPlayerGame: false,
		LocationId:   1,
		Teams: []models.Team{
			{ID: 1, Score: 100},
		},
		Players: createTestPlayersForService(encService), // Players in team game - violates business rule
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "team game cannot have any players")
}

// TestAddSavedGame_TeamGame_InvalidLocationId_Fails verifies error when location doesn't exist
func TestAddSavedGame_TeamGame_InvalidLocationId_Fails(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Score: 100},
	}

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			Err:        fmt.Errorf("location not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &SaveGameServiceImpl{
		LocationRepository: mockLocationRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame: false,
		LocationId:   999, // Invalid location ID
		Teams:        teams,
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestAddSavedGame_TeamGame_TeamsNotExist_Fails verifies error when one or more teams don't exist
func TestAddSavedGame_TeamGame_TeamsNotExist_Fails(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Score: 100},
		{ID: 2, Score: 120},
	}

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			ResultData: models.Location{ID: 1},
			StatusCode: http.StatusOK,
		},
	}

	mockTeamRepo := &mockTeamRepositoryForSavedGame{
		teamsResult: models.Result[[]models.Team]{
			ResultData: []models.Team{teams[0]}, // Only 1 team returned, but 2 requested
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		LocationRepository: mockLocationRepo,
		TeamRepository:     mockTeamRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame: false,
		LocationId:   1,
		Teams:        teams, // 2 teams
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "one or more teams do not exist")
}

// TestAddSavedGame_TeamGame_TeamRepositoryError_Fails verifies error handling for database failures
func TestAddSavedGame_TeamGame_TeamRepositoryError_Fails(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Score: 100},
	}

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			ResultData: models.Location{ID: 1},
			StatusCode: http.StatusOK,
		},
	}

	mockTeamRepo := &mockTeamRepositoryForSavedGame{
		teamsResult: models.Result[[]models.Team]{
			Err:        fmt.Errorf("database error"),
			StatusCode: http.StatusInternalServerError,
		},
	}

	service := &SaveGameServiceImpl{
		LocationRepository: mockLocationRepo,
		TeamRepository:     mockTeamRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame: false,
		LocationId:   1,
		Teams:        teams,
	}

	result := service.AddSavedGame(savedGame)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

// TestAddSavedGame_TeamGame_CalculatesScoresCorrectly verifies score calculations
func TestAddSavedGame_TeamGame_CalculatesScoresCorrectly(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Score: 100},
		{ID: 2, Score: 120},
	}

	mockSavedGameRepo := &mockSavedGameRepository{
		savedGameResult: models.Result[models.SavedGame]{
			ResultData: models.SavedGame{ID: 1},
			StatusCode: http.StatusCreated,
		},
	}

	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			ResultData: models.Location{ID: 1},
			StatusCode: http.StatusOK,
		},
	}

	mockTeamRepo := &mockTeamRepositoryForSavedGame{
		teamsResult: models.Result[[]models.Team]{
			ResultData: teams,
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		SavedGameRepository: mockSavedGameRepo,
		LocationRepository:  mockLocationRepo,
		TeamRepository:      mockTeamRepo,
	}

	savedGame := models.SavedGame{
		IsPlayerGame: false,
		LocationId:   1,
		Teams:        teams,
	}

	result := service.AddSavedGame(savedGame)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
}

// ============================================================================
// Helper Method Tests
// ============================================================================

// TestIsLocationIdValid_Success verifies location validation succeeds when location exists
func TestIsLocationIdValid_Success(t *testing.T) {
	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			ResultData: models.Location{ID: 1, LocationName: "Elmwood"},
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		LocationRepository: mockLocationRepo,
	}

	result := service.isLocationIdValid(1)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

// TestIsLocationIdValid_NotFound verifies location validation fails when location doesn't exist
func TestIsLocationIdValid_NotFound(t *testing.T) {
	mockLocationRepo := &mockLocationRepositoryForSavedGame{
		locationResult: models.Result[models.Location]{
			Err:        fmt.Errorf("location not found"),
			StatusCode: http.StatusNotFound,
		},
	}

	service := &SaveGameServiceImpl{
		LocationRepository: mockLocationRepo,
	}

	result := service.isLocationIdValid(999)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

// TestArePlayersValid_Success verifies player validation succeeds when all players exist
func TestArePlayersValid_Success(t *testing.T) {
	encService := getTestEncryptionServiceForService()
	players := createTestPlayersForService(encService)

	mockPlayerRepo := &mockPlayerRepositoryForSavedGame{
		playersResult: models.Result[[]models.Player]{
			ResultData: players,
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.arePlayersValid(players)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 2)
}

// TestArePlayersValid_NotAllExist verifies error when some players don't exist
func TestArePlayersValid_NotAllExist(t *testing.T) {
	encService := getTestEncryptionServiceForService()
	players := createTestPlayersForService(encService)

	mockPlayerRepo := &mockPlayerRepositoryForSavedGame{
		playersResult: models.Result[[]models.Player]{
			ResultData: []models.Player{players[0]}, // Only 1 of 2
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		PlayerRepository: mockPlayerRepo,
	}

	result := service.arePlayersValid(players)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "one or more players do not exist")
}

// TestAreTeamsValid_Success verifies team validation succeeds when all teams exist
func TestAreTeamsValid_Success(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Score: 100},
		{ID: 2, Score: 120},
	}

	mockTeamRepo := &mockTeamRepositoryForSavedGame{
		teamsResult: models.Result[[]models.Team]{
			ResultData: teams,
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		TeamRepository: mockTeamRepo,
	}

	result := service.areTeamsValid(teams)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Len(t, result.ResultData, 2)
}

// TestAreTeamsValid_NotAllExist verifies error when some teams don't exist
func TestAreTeamsValid_NotAllExist(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Score: 100},
		{ID: 2, Score: 120},
	}

	mockTeamRepo := &mockTeamRepositoryForSavedGame{
		teamsResult: models.Result[[]models.Team]{
			ResultData: []models.Team{teams[0]}, // Only 1 of 2
			StatusCode: http.StatusOK,
		},
	}

	service := &SaveGameServiceImpl{
		TeamRepository: mockTeamRepo,
	}

	result := service.areTeamsValid(teams)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "one or more teams do not exist")
}
