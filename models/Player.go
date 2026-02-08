package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Player struct {
	ID        			int            `db:"id"`
	CreatedAt  			time.Time      `db:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at"`
	PlayerName          string 	       `json:"-" db:"player_name"`
	PlayerNameEncrypted []byte         `db:"player_name_encrypted"`
	PlayerNameHash      []byte         `db:"player_name_hash"`
	LocationID          int            `db:"location_id"`
	TeamID              sql.NullInt32  `db:"team_id"`
	Score               int            `db:"-"`
	PlayerNameDecrypted string         `json:"-" db:"-"`
}

// /TFUYkuZM0uedGVUeqFjZw5B+kbjLidnMKigb9Wbj1E=

const (
	minLength int = 4
	maxLength int = 40
)

func (p *Player) Validate() error {
	if err := p.validatePlayerNameLength(); err != nil {
		return err
	}

	if err := p.validatePlayerNameHasTwoParts(); err != nil {
		return err
	}

	return nil
}

func (p *Player) validatePlayerNameLength() error {
	playNameLen := len(p.PlayerName)

	if playNameLen < 4 || playNameLen > 40 {
		return fmt.Errorf("player name must be between %d and %d", minLength, maxLength)
	}

	return nil
}

func (p *Player) validatePlayerNameHasTwoParts() error {
	fields := strings.Fields(p.PlayerName)

	if len(fields) != 2 {
		return fmt.Errorf("Player name must have exactly two parts: ex -> 'jane doe'")
	}

	return nil
}
