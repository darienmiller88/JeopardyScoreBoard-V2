package models

import (
	"time"
)

type Location struct {
	ID        int       `db:"id"` 
	CreatedAt time.Time `db:"created_at"` 
	UpdatedAt time.Time `db:"updated_at"`

	//The name of the ADAPT location (Pelham Bay, Lawrence, Elmwood, etc)
	LocationName string `db:"location_name"`
}