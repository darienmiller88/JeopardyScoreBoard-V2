package models

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type Location struct {
	ID        int       `db:"id"` 
	CreatedAt time.Time `db:"created_at"` 
	UpdatedAt time.Time `db:"updated_at"`

	//The name of the ADAPT location (Pelham Bay, Lawrence, Elmwood, etc)
	LocationName string `db:"location_name"`
}

func (l *Location) Validate() error{
	return validation.ValidateStruct(
		l,
		validation.Field(&l.LocationName, validation.Required),
	)
}
// This struct will represent the ADAPT community network location where a jeopardy game was played.
// type Location struct {
// 	ID        bson.ObjectID `bson:"_id,omitempty"` 
// 	CreatedAt time.Time     `bson:"created_at"` 
// 	UpdatedAt time.Time     `bson:"updated_at"`

// 	//The name of the ADAPT location (Pelham Bay, Lawrence, Elmwood, etc)
// 	LocationName string     `bson:"location_name"`

// 	//Here are all of the users that played in the game.
// 	Players     []Player `bson:"users"`
// }
