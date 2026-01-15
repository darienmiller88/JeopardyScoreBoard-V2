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

func TestSavedGameValidate_UnhappyPaths(t *testing.T) {
	tests := []struct {
		name string
		game SavedGame
	}{
		{
			name: "both winning player and team set",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Int32: 1, Valid: true},
				WinningTeamId:    sql.NullInt32{Int32: 2, Valid: true},
				WinnerPlayerName: "Jane Doe",
			},
		},
		{
			name: "neither winning player nor team set",
			game: SavedGame{
				WinningPlayerId: sql.NullInt32{Valid: false},
				WinningTeamId:   sql.NullInt32{Valid: false},
			},
		},
		{
			name: "player game missing name",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Int32: 4, Valid: true},
				WinningTeamId:    sql.NullInt32{Valid: false},
				WinnerPlayerName: "",
			},
		},
		{
			name: "player name only one part",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Int32: 4, Valid: true},
				WinningTeamId:    sql.NullInt32{Valid: false},
				WinnerPlayerName: "Jane",
			},
		},
		{
			name: "player name too many parts",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Int32: 4, Valid: true},
				WinningTeamId:    sql.NullInt32{Valid: false},
				WinnerPlayerName: "Jane Marie Doe",
			},
		},
		{
			name: "team game with winner name set",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Valid: false},
				WinningTeamId:    sql.NullInt32{Int32: 8, Valid: true},
				WinnerPlayerName: "Jane Doe",
			},
		},
	}

	for _, gameTest := range tests {
		t.Run(gameTest.name, func(t *testing.T) {
			err := gameTest.game.Validate()
			
			if err == nil {
				t.Fatalf("expected validation error but got nil")
			}
		})
	}
}
