package models

import (
	"JeopardyScoreBoardV2/encryption"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create test encryption service
func getTestEncryptionService() *encryption.EncryptionService {
	testKey := []byte("12345678901234567890123456789012") // 32 bytes for AES-256
	return encryption.NewService(testKey)
}

// Helper to create test players with encrypted names
func createTestPlayers(encService *encryption.EncryptionService) []Player {
	player1Encrypted, _ := encService.Encrypt("Player1")
	player2Encrypted, _ := encService.Encrypt("Player2")
	player3Encrypted, _ := encService.Encrypt("Player3")

	return []Player{
		{
			ID:                  1,
			PlayerName:          "Player1",
			PlayerNameEncrypted: player1Encrypted,
			PlayerNameHash:      encryption.NameHash("Player1"),
			Score:               50,
		},
		{
			ID:                  2,
			PlayerName:          "Player2",
			PlayerNameEncrypted: player2Encrypted,
			PlayerNameHash:      encryption.NameHash("Player2"),
			Score:               60,
		},
		{
			ID:                  3,
			PlayerName:          "Player3",
			PlayerNameEncrypted: player3Encrypted,
			PlayerNameHash:      encryption.NameHash("Player3"),
			Score:               40,
		},
	}
}

// ============================================================================
// Validate() Tests - Table Driven
// ============================================================================

func TestSavedGame_Validate(t *testing.T) {
	encService := getTestEncryptionService()

	tests := []struct {
		name          string
		savedGame     SavedGame
		expectedError string
	}{
		{
			name: "Player game valid - happy",
			savedGame: SavedGame{
				IsPlayerGame:      true,
				Players:           createTestPlayers(encService)[:2],
				Teams:             []Team{},
				EncryptionService: encService,
			},
			expectedError: "",
		},
		{
			name: "Team game valid - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Players:      []Player{},
				Teams: []Team{
					{ID: 1, Score: 100},
					{ID: 2, Score: 120},
				},
				EncryptionService: encService,
			},
			expectedError: "",
		},
		{
			name: "Player game with teams - unhappy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players:      createTestPlayers(encService)[:1],
				Teams: []Team{
					{ID: 1, Score: 100},
				},
				EncryptionService: encService,
			},
			expectedError: "a player game cannot have any teams added",
		},
		{
			name: "Team game with players - unhappy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Players:      createTestPlayers(encService)[:1],
				Teams: []Team{
					{ID: 1, Score: 100},
				},
				EncryptionService: encService,
			},
			expectedError: "a team game cannot have any players added",
		},
		{
			name: "Player game empty players - unhappy",
			savedGame: SavedGame{
				IsPlayerGame:      true,
				Players:           []Player{},
				Teams:             []Team{},
				EncryptionService: encService,
			},
			expectedError: "players cannot be empty when a player game is being played",
		},
		{
			name: "Team game empty teams - unhappy",
			savedGame: SavedGame{
				IsPlayerGame:      false,
				Players:           []Player{},
				Teams:             []Team{},
				EncryptionService: encService,
			},
			expectedError: "teams cannot be empty when a team game is being played",
		},
		{
			name: "Player game nil players - unhappy",
			savedGame: SavedGame{
				IsPlayerGame:      true,
				Players:           nil,
				EncryptionService: encService,
			},
			expectedError: "players cannot be empty when a player game is being played",
		},
		{
			name: "Team game nil teams - unhappy",
			savedGame: SavedGame{
				IsPlayerGame:      false,
				Teams:             nil,
				EncryptionService: encService,
			},
			expectedError: "teams cannot be empty when a team game is being played",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.savedGame.Validate()

			if tt.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			}
		})
	}
}

// ============================================================================
// CalculateTotalPoints() Tests - Table Driven
// ============================================================================

func TestSavedGame_CalculateTotalPoints(t *testing.T) {
	encService := getTestEncryptionService()

	tests := []struct {
		name          string
		savedGame     SavedGame
		expectedTotal int
	}{
		{
			name: "Player game multiple players - happy",
			savedGame: SavedGame{
				IsPlayerGame:      true,
				Players:           createTestPlayers(encService),
				EncryptionService: encService,
			},
			expectedTotal: 150,
		},
		{
			name: "Team game multiple teams - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, Score: 100},
					{ID: 2, Score: 120},
					{ID: 3, Score: 80},
				},
				EncryptionService: encService,
			},
			expectedTotal: 300,
		},
		{
			name: "Team game multiple teams and negative scores - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, Score: -100},
					{ID: 2, Score: 120},
					{ID: 3, Score: -80},
				},
				EncryptionService: encService,
			},
			expectedTotal: -60,
		},
		{
			name: "Player game empty - happy",
			savedGame: SavedGame{
				IsPlayerGame:      true,
				Players:           []Player{},
				EncryptionService: encService,
			},
			expectedTotal: 0,
		},
		{
			name: "Team game empty - happy",
			savedGame: SavedGame{
				IsPlayerGame:      false,
				Teams:             []Team{},
				EncryptionService: encService,
			},
			expectedTotal: 0,
		},
		{
			name: "Player game single player - happy",
			savedGame: SavedGame{
				IsPlayerGame:      true,
				Players:           createTestPlayers(encService)[:1],
				EncryptionService: encService,
			},
			expectedTotal: 50,
		},
		{
			name: "Team game single team - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, Score: 150},
				},
				EncryptionService: encService,
			},
			expectedTotal: 150,
		},
		{
			name: "Player game zero scores - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, Score: 0, PlayerNameEncrypted: []byte("encrypted1")},
					{ID: 2, Score: 0, PlayerNameEncrypted: []byte("encrypted2")},
				},
				EncryptionService: encService,
			},
			expectedTotal: 0,
		},
		{
			name: "Resets existing total - happy",
			savedGame: SavedGame{
				IsPlayerGame:      true,
				TotalPoints:       999,
				Players:           createTestPlayers(encService)[:2],
				EncryptionService: encService,
			},
			expectedTotal: 110,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.savedGame.CalculateTotalPoints()
			assert.Equal(t, tt.expectedTotal, tt.savedGame.TotalPoints)
		})
	}
}

// ============================================================================
// CalculateAveragePoints() Tests - Table Driven
// ============================================================================

func TestSavedGame_CalculateAveragePoints(t *testing.T) {
	encService := getTestEncryptionService()

	tests := []struct {
		name            string
		savedGame       SavedGame
		expectedAverage float64
		expectedTotal   int
	}{
		{
			name: "Player game even division - happy",
			savedGame: SavedGame{
				IsPlayerGame:      true,
				Players:           createTestPlayers(encService),
				EncryptionService: encService,
			},
			expectedAverage: 50.0,
			expectedTotal:   150,
		},
		{
			name: "Team game even division - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, Score: 100},
					{ID: 2, Score: 120},
				},
				EncryptionService: encService,
			},
			expectedAverage: 110.0,
			expectedTotal:   220,
		},
		{
			name: "Player game single player - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{
						ID:                  1,
						Score:               75,
						PlayerNameEncrypted: []byte("encrypted"),
					},
				},
				EncryptionService: encService,
			},
			expectedAverage: 75.0,
			expectedTotal:   75,
		},
		{
			name: "Player game non-even division - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, Score: 50, PlayerNameEncrypted: []byte("enc1")},
					{ID: 2, Score: 60, PlayerNameEncrypted: []byte("enc2")},
					{ID: 3, Score: 45, PlayerNameEncrypted: []byte("enc3")},
				},
				EncryptionService: encService,
			},
			expectedAverage: 51.666666,
			expectedTotal:   155,
		},
		{
			name: "Team game non-even division - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, Score: 100},
					{ID: 2, Score: 110},
					{ID: 3, Score: 105},
				},
				EncryptionService: encService,
			},
			expectedAverage: 105.0,
			expectedTotal:   315,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.savedGame.CalculateAveragePoints()
			assert.InDelta(t, tt.expectedAverage, tt.savedGame.AveragePoints, 0.0001)
			assert.Equal(t, tt.expectedTotal, tt.savedGame.TotalPoints)
		})
	}
}

// ============================================================================
// CalculateWinner() Tests - Table Driven
// ============================================================================

func TestSavedGame_CalculateWinner(t *testing.T) {
	encService := getTestEncryptionService()

	tests := []struct {
		name                     string
		savedGame                SavedGame
		expectedPlayerName       string
		expectedPlayerId         int32
		expectedTeamId           int32
		shouldHaveEncryptedName  bool
		shouldHavePlayerId       bool
		shouldHaveTeamId         bool
	}{
		{
			name: "Player game single player - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{
						ID:                  1,
						PlayerName:          "Player1",
						PlayerNameEncrypted: mustEncrypt(encService, "Player1"),
						Score:               75,
					},
				},
				EncryptionService: encService,
			},
			expectedPlayerName:      "Player1",
			expectedPlayerId:        1,
			shouldHaveEncryptedName: true,
			shouldHavePlayerId:      true,
			shouldHaveTeamId:        false,
		},
		{
			name: "Team game single team - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, Score: 150},
				},
				EncryptionService: encService,
			},
			expectedTeamId:          1,
			shouldHaveEncryptedName: false,
			shouldHavePlayerId:      false,
			shouldHaveTeamId:        true,
		},
		{
			name: "Player game multiple players - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", PlayerNameEncrypted: mustEncrypt(encService, "Player1"), Score: 50},
					{ID: 2, PlayerName: "Player2", PlayerNameEncrypted: mustEncrypt(encService, "Player2"), Score: 75},
					{ID: 3, PlayerName: "Player3", PlayerNameEncrypted: mustEncrypt(encService, "Player3"), Score: 60},
				},
				EncryptionService: encService,
			},
			expectedPlayerName:      "Player2",
			expectedPlayerId:        2,
			shouldHaveEncryptedName: true,
			shouldHavePlayerId:      true,
			shouldHaveTeamId:        false,
		},
		{
			name: "Team game multiple teams - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, Score: 100},
					{ID: 2, Score: 150},
					{ID: 3, Score: 120},
				},
				EncryptionService: encService,
			},
			expectedTeamId:          2,
			shouldHaveEncryptedName: false,
			shouldHavePlayerId:      false,
			shouldHaveTeamId:        true,
		},
		{
			name: "Player game tied scores - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", PlayerNameEncrypted: mustEncrypt(encService, "Player1"), Score: 75},
					{ID: 2, PlayerName: "Player2", PlayerNameEncrypted: mustEncrypt(encService, "Player2"), Score: 75},
					{ID: 3, PlayerName: "Player3", PlayerNameEncrypted: mustEncrypt(encService, "Player3"), Score: 60},
				},
				EncryptionService: encService,
			},
			expectedPlayerName:      "Player1",
			expectedPlayerId:        1,
			shouldHaveEncryptedName: true,
			shouldHavePlayerId:      true,
			shouldHaveTeamId:        false,
		},
		{
			name: "Team game tied scores - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, Score: 150},
					{ID: 2, Score: 150},
					{ID: 3, Score: 120},
				},
				EncryptionService: encService,
			},
			expectedTeamId:          1,
			shouldHaveEncryptedName: false,
			shouldHavePlayerId:      false,
			shouldHaveTeamId:        true,
		},
		{
			name: "Player game two players - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", PlayerNameEncrypted: mustEncrypt(encService, "Player1"), Score: 50},
					{ID: 2, PlayerName: "Player2", PlayerNameEncrypted: mustEncrypt(encService, "Player2"), Score: 75},
				},
				EncryptionService: encService,
			},
			expectedPlayerName:      "Player2",
			expectedPlayerId:        2,
			shouldHaveEncryptedName: true,
			shouldHavePlayerId:      true,
			shouldHaveTeamId:        false,
		},
		{
			name: "Team game two teams - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, Score: 100},
					{ID: 2, Score: 150},
				},
				EncryptionService: encService,
			},
			expectedTeamId:          2,
			shouldHaveEncryptedName: false,
			shouldHavePlayerId:      false,
			shouldHaveTeamId:        true,
		},
		{
			name: "Player game winner is first - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", PlayerNameEncrypted: mustEncrypt(encService, "Player1"), Score: 100},
					{ID: 2, PlayerName: "Player2", PlayerNameEncrypted: mustEncrypt(encService, "Player2"), Score: 75},
					{ID: 3, PlayerName: "Player3", PlayerNameEncrypted: mustEncrypt(encService, "Player3"), Score: 60},
				},
				EncryptionService: encService,
			},
			expectedPlayerName:      "Player1",
			expectedPlayerId:        1,
			shouldHaveEncryptedName: true,
			shouldHavePlayerId:      true,
			shouldHaveTeamId:        false,
		},
		{
			name: "Player game winner is last - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", PlayerNameEncrypted: mustEncrypt(encService, "Player1"), Score: 50},
					{ID: 2, PlayerName: "Player2", PlayerNameEncrypted: mustEncrypt(encService, "Player2"), Score: 60},
					{ID: 3, PlayerName: "Player3", PlayerNameEncrypted: mustEncrypt(encService, "Player3"), Score: 100},
				},
				EncryptionService: encService,
			},
			expectedPlayerName:      "Player3",
			expectedPlayerId:        3,
			shouldHaveEncryptedName: true,
			shouldHavePlayerId:      true,
			shouldHaveTeamId:        false,
		},
		{
			name: "Player game zero scores - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", PlayerNameEncrypted: mustEncrypt(encService, "Player1"), Score: 0},
					{ID: 2, PlayerName: "Player2", PlayerNameEncrypted: mustEncrypt(encService, "Player2"), Score: 0},
				},
				EncryptionService: encService,
			},
			expectedPlayerName:      "Player1",
			expectedPlayerId:        1,
			shouldHaveEncryptedName: true,
			shouldHavePlayerId:      true,
			shouldHaveTeamId:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.savedGame.CalculateWinner()
			require.NoError(t, err)

			if tt.shouldHaveEncryptedName {
				assert.NotEmpty(t, tt.savedGame.WinningPlayerNameEncrypted)
				// Decrypt and verify
				decrypted, err := encService.Decrypt(tt.savedGame.WinningPlayerNameEncrypted)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPlayerName, decrypted)
			}

			if tt.shouldHavePlayerId {
				assert.True(t, tt.savedGame.WinningPlayerId.Valid)
				assert.Equal(t, tt.expectedPlayerId, tt.savedGame.WinningPlayerId.Int32)
			}

			if tt.shouldHaveTeamId {
				assert.True(t, tt.savedGame.WinningTeamId.Valid)
				assert.Equal(t, tt.expectedTeamId, tt.savedGame.WinningTeamId.Int32)
			}
		})
	}
}

// Helper function to encrypt or panic (for test setup only)
func mustEncrypt(encService *encryption.EncryptionService, plaintext string) []byte {
	encrypted, err := encService.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}
	return encrypted
}

// ============================================================================
// Helper Function Tests - Table Driven
// ============================================================================

func TestNewNullInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int32
	}{
		{name: "Positive number - happy", input: 42, expected: 42},
		{name: "Zero - happy", input: 0, expected: 0},
		{name: "Negative number - happy", input: -5, expected: -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newNullInt32(tt.input)
			assert.True(t, result.Valid)
			assert.Equal(t, tt.expected, result.Int32)
		})
	}
}

func TestNewNullString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Normal string - happy", input: "Test String", expected: "Test String"},
		{name: "Empty string - happy", input: "", expected: ""},
		{name: "Special characters - happy", input: "Test!@#$%^&*()", expected: "Test!@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newNullString(tt.input)
			assert.True(t, result.Valid)
			assert.Equal(t, tt.expected, result.String)
		})
	}
}