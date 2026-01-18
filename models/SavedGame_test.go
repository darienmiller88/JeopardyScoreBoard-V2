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
				WinningPlayerName: sql.NullString{String: "vicky Miller", Valid: true},
				Players: []Player{
					{ PlayerName: "Darien Miller" },
					{ PlayerName: "vicky Miller" },
				},
			},
		},
		{
			name: "valid team based game",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Valid: false},
				WinningTeamId:    sql.NullInt32{Int32: 3, Valid: true},
				WinningPlayerName: sql.NullString{Valid: false},
				Teams: []Team{
					{ LocationID: 1 },
					{ LocationID: 2 },
				},
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
			name: "player game missing winning player name",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Int32: 4, Valid: true},
				WinningTeamId:    sql.NullInt32{Valid: false},
				WinningPlayerName: sql.NullString{Valid: false},
			},
		},
		{
			name: "winning player name only one part",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Int32: 4, Valid: true},
				WinningTeamId:    sql.NullInt32{Valid: false},
				WinningPlayerName: sql.NullString{String:"Jane", Valid: true},
			},
		},
		{
			name: "winning player name too many parts",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Int32: 4, Valid: true},
				WinningTeamId:    sql.NullInt32{Valid: false},
				WinningPlayerName: sql.NullString{String: "Jane Marie Doe", Valid: true},
			},
		},
		{
			name: "team game with winner player name",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Valid: false},
				WinningTeamId:    sql.NullInt32{Int32: 8, Valid: true},
				WinningPlayerName: sql.NullString{String: "Jane Doe", Valid: true},
			},
		},
		{
			name: "team game with no teams",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Valid: false},
				WinningTeamId:    sql.NullInt32{Valid: true},
			},
		},
		{
			name: "player game with no players",
			game: SavedGame{
				WinningPlayerId:  sql.NullInt32{Valid: true},
				WinningTeamId:    sql.NullInt32{Valid: false},
			},
		},
		{
			name: "player game with no winning player",
			game: SavedGame{
				WinningPlayerId:   sql.NullInt32{Valid: true},
				WinningTeamId:     sql.NullInt32{Valid: false},
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
