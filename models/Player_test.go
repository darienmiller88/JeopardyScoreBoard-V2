package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test for getting all 8 locations
//EXPECTED PAYLOAD:       All 8 locations
//EXPECTED STATUS CODE:   200
//EXPECTED ERROR MESSAGE: nil   
func TestPlayerModelValidation_Ok(t *testing.T) {
	player := Player{
		PlayerName: "Test Name",
	}

	err := player.Validate()

	assert.Equal(t, nil, err)
}

func TestPlayerModelValidation_NameTooShort(t *testing.T) {
	
}

func TestPlayerModelValidation_NameTooLong(t *testing.T) {
	
}


func TestPlayerModelValidation_NameHasOnePart(t *testing.T) {
	
}

func TestPlayerModelValidation_NameHasThreeParts(t *testing.T) {
	
}