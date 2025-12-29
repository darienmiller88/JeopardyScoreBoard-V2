package models

import (
	"database/sql"
	"time"
)

type Player struct{ 
	ID         int           `db:"id"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`
	PlayerName string        `db:"player_name"`
	LocationID int           `db:"location_id"`
	TeamID     sql.NullInt32 `db:"team_id"`
}

// func (p *Player) Validate() error{
// 	return validation.ValidateStruct(
// 		p,
// 		validation.Field(&p.Name, validation.Length(3, 31)),
// 	)
// }