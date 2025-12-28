package models

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type Player struct{ 
	Score int    `bson:"score"` 
	Name  string `bson:"name"`
}

func (p *Player) Validate() error{
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Name, validation.Length(3, 31)),
	)
}