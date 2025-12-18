package models

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// This struct will represent the ADAPT community network location where a jeopardy game was played.
type Location struct {
	ID        bson.ObjectID `bson:"_id,omitempty"` 
	CreatedAt time.Time     `bson:"created_at"` 
	UpdatedAt time.Time     `bson:"updated_at"`

	//The name of the ADAPT location (Pelham Bay, Lawrence, Elmwood, etc)
	LocationName string     `bson:"location_name"`

	//Here are all of the users that played in the game.
	Players     []PlayerCard `bson:"users"`
}

func (l *Location) InitCreatedAtAndUpdatedAt(){
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()

	//If this field is not initialized, it is interpreted as "null" by mongoDB, and not an empty array.
	l.Players = []PlayerCard{}
}

func (l *Location) Validate() error{
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Players, validation.Required, validation.Length(1, 0)),
		validation.Field(&l.LocationName, validation.Required),
	)
}