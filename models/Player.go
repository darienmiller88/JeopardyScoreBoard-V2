package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type Player struct{ 
	ID         int           `db:"id"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`
	PlayerName string        `db:"player_name"`
	LocationID int           `db:"location_id"`
	TeamID     sql.NullInt32 `db:"team_id"`
}

func (p *Player) Validate() error{
	return validation.ValidateStruct(
		p,
		validation.Field(&p.PlayerName, validation.Required, validation.Length(5, 40), validation.By(p.validatePlayerNameHasTwoParts)),
	)
}

func (p *Player) validatePlayerNameHasTwoParts(field interface{}) error{
	playerName, ok := field.(string)

	if !ok{
		return fmt.Errorf("could not parse %T into object", field)
	}
	
	fields := strings.Fields(playerName)

	if len(fields) != 2 {
		return fmt.Errorf("Player name must have exactly two parts: ex -> 'jane doe'")
	}

	return nil
}