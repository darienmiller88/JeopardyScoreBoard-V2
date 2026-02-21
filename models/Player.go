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
	PlayerNameEncrypted []byte         `db:"player_name_encrypted"`
	PlayerNameHash      []byte         `db:"player_name_hash"`
	LocationID          int            `db:"location_id"`
	TeamID              sql.NullInt32  `db:"team_id"`
	
	// Non DB fields
	Score               int            `db:"-"`
	PlayerName          string 	       
	FirstName           string
	LastName            string
	// PlayerNameDecrypted string         `json:"-" db:"-"`
	PlayerNameAbbrev    string         ``
}

const (
	minLength int = 4
	maxLength int = 40

	// min and max Length for first and last name
	minNameLength = 4
	maxNameLength = 20
	
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

func (p *Player) SetPlayerName(firstName string, lastName string){
	p.PlayerName = firstName + " " + lastName
}

func (p *Player) validatePlayerNameLength() error {
	firstNameLen := len(p.FirstName)
	lastNameLen := len(p.LastName)

	if (firstNameLen < minNameLength || firstNameLen > maxNameLength) {
		return fmt.Errorf("First name must be between %d and %d", minNameLength, maxNameLength)
	}

	if (lastNameLen < minNameLength || lastNameLen > maxNameLength) {
		return fmt.Errorf("Last name must be between %d and %d", minNameLength, maxNameLength)
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
