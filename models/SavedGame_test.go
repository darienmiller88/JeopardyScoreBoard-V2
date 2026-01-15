package models

import (
	"database/sql"
	"testing"
)

func TestSavedGameValidate_HappyPaths(t *testing.T) {
	tests := []struct {
		name string
		game SavedGame
	}{
		{
			name: "valid player based game",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Int32: 5, Valid: true},
				WinningTeamId:    sql.NullInt32{Valid: false},
				WinnerPlayerName: "Jane Doe",
			},
		},
		{
			name: "valid team based game",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Valid: false},
				WinningTeamId:    sql.NullInt32{Int32: 3, Valid: true},
				WinnerPlayerName: "",
			},
		},
	}

	for _, gameTest := range tests {
		t.Run(gameTest.name, func(t *testing.T) {
			err := gameTest.game.Validate()
			
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
