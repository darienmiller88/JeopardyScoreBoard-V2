package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test for a valid player name
//EXPECTED ERROR MESSAGE: nil
func TestPlayerModelValidation_Ok(t *testing.T) {
	player := Player{
		PlayerName: "Test Name",
	}

	err := player.Validate()

	assert.Equal(t, nil, err)
}

//Test for a valid player name
//EXPECTED ERROR MESSAGE: nil 
func TestPlayerModelValidation_NameTooShort(t *testing.T) {
	player := Player{
		PlayerName: "Te",
	}

	err := player.Validate()

	assert.Contains(t, err.Error(), fmt.Sprintf("player name must be between %d and %d", minLength, maxLength))
}

func TestPlayerModelValidation_NameTooLong(t *testing.T) {
	player := Player{
		PlayerName: "Tdedxchufedxiucfixmuebfxdbxzwubxzbinwjbdhzsedbxywhzsnuebyfdxzwse",
	}

	err := player.Validate()

	assert.Contains(t, err.Error(), fmt.Sprintf("player name must be between %d and %d", minLength, maxLength))
}


func TestPlayerModelValidation_NameHasOnePart(t *testing.T) {
	player := Player{
		PlayerName: "ThisnameisOnePart",
	}

	err := player.Validate()

	assert.Contains(t, err.Error(), "Player name must have exactly two parts: ex -> 'jane doe'")
}

func TestPlayerModelValidation_NameHasMoreThanTwoParts(t *testing.T) {
	player := Player{
		PlayerName: "   This name is multiple Parts   ",
	}

	err := player.Validate()

	assert.Contains(t, err.Error(), "Player name must have exactly two parts: ex -> 'jane doe'")
}