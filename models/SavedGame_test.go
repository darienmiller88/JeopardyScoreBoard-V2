package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Validate() Tests - Table Driven
// ============================================================================

func TestSavedGame_Validate(t *testing.T) {
	tests := []struct {
		name          string
		savedGame     SavedGame
		expectedError string
	}{
		{
			name: "Player game valid - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
					{ID: 2, PlayerName: "Player2", Score: 60},
				},
				Teams: []Team{},
			},
			expectedError: "",
		},
		{
			name: "Team game valid - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Players:      []Player{},
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 100},
					{ID: 2, TeamName: "Team2", Score: 120},
				},
			},
			expectedError: "",
		},
		{
			name: "Player game with teams - unhappy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
				},
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 100},
				},
			},
			expectedError: "a player game cannot have any teams added",
		},
		{
			name: "Team game with players - unhappy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
				},
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 100},
				},
			},
			expectedError: "a team game cannot have any players added",
		},
		{
			name: "Player game empty players - unhappy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players:      []Player{},
				Teams:        []Team{},
			},
			expectedError: "players cannot be empty when winning player id is supplied",
		},
		{
			name: "Team game empty teams - unhappy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Players:      []Player{},
				Teams:        []Team{},
			},
			expectedError: "teams cannot be empty when winning team id is supplied",
		},
		{
			name: "Player game nil players - unhappy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players:      nil,
			},
			expectedError: "players cannot be empty when winning player id is supplied",
		},
		{
			name: "Team game nil teams - unhappy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams:        nil,
			},
			expectedError: "teams cannot be empty when winning team id is supplied",
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
	tests := []struct {
		name          string
		savedGame     SavedGame
		expectedTotal int
	}{
		{
			name: "Player game multiple players - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
					{ID: 2, PlayerName: "Player2", Score: 60},
					{ID: 3, PlayerName: "Player3", Score: 40},
				},
			},
			expectedTotal: 150,
		},
		{
			name: "Team game multiple teams - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 100},
					{ID: 2, TeamName: "Team2", Score: 120},
					{ID: 3, TeamName: "Team3", Score: 80},
				},
			},
			expectedTotal: 300,
		},
		{
			name: "Player game empty - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players:      []Player{},
			},
			expectedTotal: 0,
		},
		{
			name: "Team game empty - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams:        []Team{},
			},
			expectedTotal: 0,
		},
		{
			name: "Player game single player - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 75},
				},
			},
			expectedTotal: 75,
		},
		{
			name: "Team game single team - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 150},
				},
			},
			expectedTotal: 150,
		},
		{
			name: "Player game zero scores - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 0},
					{ID: 2, PlayerName: "Player2", Score: 0},
				},
			},
			expectedTotal: 0,
		},
		{
			name: "Resets existing total - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				TotalPoints:  999,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
					{ID: 2, PlayerName: "Player2", Score: 60},
				},
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
	tests := []struct {
		name            string
		savedGame       SavedGame
		expectedAverage float64
		expectedTotal   int
	}{
		{
			name: "Player game even division - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
					{ID: 2, PlayerName: "Player2", Score: 60},
					{ID: 3, PlayerName: "Player3", Score: 40},
				},
			},
			expectedAverage: 50.0,
			expectedTotal:   150,
		},
		{
			name: "Team game even division - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 100},
					{ID: 2, TeamName: "Team2", Score: 120},
				},
			},
			expectedAverage: 110.0,
			expectedTotal:   220,
		},
		{
			name: "Player game single player - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 75},
				},
			},
			expectedAverage: 75.0,
			expectedTotal:   75,
		},
		{
			name: "Player game non-even division - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
					{ID: 2, PlayerName: "Player2", Score: 60},
					{ID: 3, PlayerName: "Player3", Score: 45},
				},
			},
			expectedAverage: 51.666666,
			expectedTotal:   155,
		},
		{
			name: "Team game non-even division - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 100},
					{ID: 2, TeamName: "Team2", Score: 110},
					{ID: 3, TeamName: "Team3", Score: 105},
				},
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
	tests := []struct {
		name                  string
		savedGame             SavedGame
		expectedPlayerName    string
		expectedPlayerId      int32
		expectedTeamId        int32
		shouldHavePlayerName  bool
		shouldHavePlayerId    bool
		shouldHaveTeamId      bool
	}{
		{
			name: "Player game single player - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 75},
				},
			},
			expectedPlayerName:   "Player1",
			expectedPlayerId:     1,
			shouldHavePlayerName: true,
			shouldHavePlayerId:   true,
			shouldHaveTeamId:     false,
		},
		{
			name: "Team game single team - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 150},
				},
			},
			expectedTeamId:       1,
			shouldHavePlayerName: false,
			shouldHavePlayerId:   false,
			shouldHaveTeamId:     true,
		},
		{
			name: "Player game multiple players - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
					{ID: 2, PlayerName: "Player2", Score: 75},
					{ID: 3, PlayerName: "Player3", Score: 60},
				},
			},
			expectedPlayerName:   "Player2",
			expectedPlayerId:     2,
			shouldHavePlayerName: true,
			shouldHavePlayerId:   true,
			shouldHaveTeamId:     false,
		},
		{
			name: "Team game multiple teams - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 100},
					{ID: 2, TeamName: "Team2", Score: 150},
					{ID: 3, TeamName: "Team3", Score: 120},
				},
			},
			expectedTeamId:       2,
			shouldHavePlayerName: false,
			shouldHavePlayerId:   false,
			shouldHaveTeamId:     true,
		},
		{
			name: "Player game tied scores - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 75},
					{ID: 2, PlayerName: "Player2", Score: 75},
					{ID: 3, PlayerName: "Player3", Score: 60},
				},
			},
			expectedPlayerName:   "Player1",
			expectedPlayerId:     1,
			shouldHavePlayerName: true,
			shouldHavePlayerId:   true,
			shouldHaveTeamId:     false,
		},
		{
			name: "Team game tied scores - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 150},
					{ID: 2, TeamName: "Team2", Score: 150},
					{ID: 3, TeamName: "Team3", Score: 120},
				},
			},
			expectedTeamId:       1,
			shouldHavePlayerName: false,
			shouldHavePlayerId:   false,
			shouldHaveTeamId:     true,
		},
		{
			name: "Player game two players - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
					{ID: 2, PlayerName: "Player2", Score: 75},
				},
			},
			expectedPlayerName:   "Player2",
			expectedPlayerId:     2,
			shouldHavePlayerName: true,
			shouldHavePlayerId:   true,
			shouldHaveTeamId:     false,
		},
		{
			name: "Team game two teams - happy",
			savedGame: SavedGame{
				IsPlayerGame: false,
				Teams: []Team{
					{ID: 1, TeamName: "Team1", Score: 100},
					{ID: 2, TeamName: "Team2", Score: 150},
				},
			},
			expectedTeamId:       2,
			shouldHavePlayerName: false,
			shouldHavePlayerId:   false,
			shouldHaveTeamId:     true,
		},
		{
			name: "Player game winner is first - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 100},
					{ID: 2, PlayerName: "Player2", Score: 75},
					{ID: 3, PlayerName: "Player3", Score: 60},
				},
			},
			expectedPlayerName:   "Player1",
			expectedPlayerId:     1,
			shouldHavePlayerName: true,
			shouldHavePlayerId:   true,
			shouldHaveTeamId:     false,
		},
		{
			name: "Player game winner is last - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 50},
					{ID: 2, PlayerName: "Player2", Score: 60},
					{ID: 3, PlayerName: "Player3", Score: 100},
				},
			},
			expectedPlayerName:   "Player3",
			expectedPlayerId:     3,
			shouldHavePlayerName: true,
			shouldHavePlayerId:   true,
			shouldHaveTeamId:     false,
		},
		{
			name: "Player game zero scores - happy",
			savedGame: SavedGame{
				IsPlayerGame: true,
				Players: []Player{
					{ID: 1, PlayerName: "Player1", Score: 0},
					{ID: 2, PlayerName: "Player2", Score: 0},
				},
			},
			expectedPlayerName:   "Player1",
			expectedPlayerId:     1,
			shouldHavePlayerName: true,
			shouldHavePlayerId:   true,
			shouldHaveTeamId:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.savedGame.CalculateWinner()

			if tt.shouldHavePlayerName {
				assert.True(t, tt.savedGame.WinningPlayerName.Valid)
				assert.Equal(t, tt.expectedPlayerName, tt.savedGame.WinningPlayerName.String)
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
