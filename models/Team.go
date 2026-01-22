package models

import (
	"time"
)

type Team struct{
	ID         int        `db:"id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	LocationID int        `db:"location_id"`
	Score      int        `db:"-"` 
	Players    []Player   `db:"-"`
	TeamName   string     `db:"-"`
}