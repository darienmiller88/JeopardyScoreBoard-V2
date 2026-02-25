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
	PlayerNameAbbrev    string         ``
}

const (
	// min and max Length for first and last name
	minNameLength = 2
	maxNameLength = 20
)

func (p *Player) Validate() error {
	if err := p.validateFirstAndLastNameLength(); err != nil {
		return err
	}

	if err := p.validateFirstAndLastNameHasOnePart(); err != nil {
		return err
	}

	return nil
}

//Combines the first name and 
func (p *Player) SetPlayerName(firstName string, lastName string){
	p.PlayerName = firstName + " " + lastName
}

func (p *Player) validateFirstAndLastNameLength() error {
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

func (p *Player) validateFirstAndLastNameHasOnePart() error {
	firstNameFields := strings.Fields(p.FirstName)

	if len(firstNameFields) != 1{
		return fmt.Errorf("First name must have exactly 1 part: ex -> 'jane'")
	}

	lastNameFields := strings.Fields(p.LastName)

	if len(lastNameFields) != 1{
		return fmt.Errorf("Last name must have exactly 1 part: ex -> 'Doe'")
	}

	return nil
}
