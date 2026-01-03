package models

import (
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

	assert.Equal(t, "edssxw", err)
}

func TestPlayerModelValidation_NameTooLong(t *testing.T) {
	
}


func TestPlayerModelValidation_NameHasOnePart(t *testing.T) {
	
}

func TestPlayerModelValidation_NameHasThreeParts(t *testing.T) {
	
}