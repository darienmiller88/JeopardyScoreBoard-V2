package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type SavedGame struct {
	ID               int           `db:"id"` 
	CreatedAt        time.Time     `db:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at"`
	TotalPoints      int 		   `db:"total_score"`
	AveragePoints    float64       `db:"average_score"`
	LocationId       int           `db:"location_id"`
	WinnerPlayerName string        `db:"winner_player_name"`
	WinningTeamId    sql.NullInt32 `db:"winning_team_id"`
	WinningPlayerId  sql.NullInt32 `db:"player_id"`
}

func (s *SavedGame) Validate() error{
	if err := s.validateGameType(); err != nil {
		return err
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
		return errors.New("saved game must have either a winning team id and winning player id")
	}

	//If there is a player id, it must also be accompanied with a valid winning player name
	if s.WinningPlayerId.Valid && s.validateWinningPlayerNameHasTwoParts() != nil {
		return errors.New("saved game cannot have both a winning team id and winning player id")
	}

	//If the saved game is a team game, there can be no winning player name 
	if s.WinningTeamId.Valid && s.WinnerPlayerName != "" {
		return errors.New("team game cannot have a winning player, only winning team id")
	}

	return nil
}

func (s *SavedGame) validateWinningPlayerNameHasTwoParts() error{	
	fields := strings.Fields(s.WinnerPlayerName)

	if len(fields) != 2 {
		return errors.New("Player name must have exactly two parts: ex -> 'jane doe'")
	}

	return nil
}