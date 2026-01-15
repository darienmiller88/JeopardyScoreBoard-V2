package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
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
	return validation.ValidateStruct(
		s,

		//Order matters here! Validate the location name first to ensure that the Location field is not nil

	)
}

func (s *SavedGame) validateGameType() error{
	if s.WinningPlayerId.Valid && s.WinningTeamId.Valid {
		return errors.New("Teams")
	}

	return nil
}

func (s *SavedGame) validateWinningPlayerName() error{
	return nil
}

func (s *SavedGame) validatePoints() error{
	return nil
}