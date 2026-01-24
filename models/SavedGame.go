package models

import (
	"database/sql"
	"errors"
	"time"
)

type SavedGame struct {
	ID                int            `db:"id"` 
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
	TotalPoints       int 		     `db:"total_score"`
	AveragePoints     float64        `db:"average_score"`
	LocationId        int            `db:"location_id"`
	WinningPlayerName sql.NullString `db:"winning_player_name"`
	WinningTeamId     sql.NullInt32  `db:"winning_team_id"`
	WinningPlayerId   sql.NullInt32  `db:"winning_player_id"`
	Players           []Player       `json:"players" db:"-"`
	Teams             []Team         `json:"teams"   db:"-"`
}

func (s *SavedGame) Validate() error{
	if err := s.validateGameType(); err != nil {
		return err
	}

	if len(s.Players) == 0 && s.WinningPlayerId.Valid{
		return errors.New("players cannot be empty when winning player id is supplied")
	}

	if len(s.Teams) == 0 && s.WinningTeamId.Valid{
		return errors.New("teams cannot be empty when winning team id is supplied")
	}

	return nil
}

func (s *SavedGame) validateGameType() error{
	//Enforce mutual exclusivity of saved games either being a team game or a player based game
	if s.WinningPlayerId.Valid && s.WinningTeamId.Valid {
		return errors.New("saved game cannot have both a winning team id and winning player id")
	}

	//Of course, both the winnning player id and winning team id can't be null
	if !s.WinningPlayerId.Valid && !s.WinningTeamId.Valid {
		return errors.New("saved game cannot have both a null a winning team id and winning player id")
	}

	//If the saved game is a team game, there can be no winning player name 
	if s.WinningTeamId.Valid && s.WinningPlayerName.Valid {
		return errors.New("team game cannot have a winning player, only winning team id")
	}

	return nil
}