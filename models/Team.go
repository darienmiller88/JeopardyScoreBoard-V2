package models

import (
	"fmt"
	"strings"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type Team struct{
	Score     int      `bson:"score"` 
	TeamName  string   `bson:"team_name"`
	Players   []string `bson:"players"`
	
}

func (t *Team) Validate() error{
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Players, validation.Required, validation.Length(2, 0), validation.By(t.checkDuplicateTeamPlayers)),
	)
}

//Check for duplicates in the list of users sent by the client.
func (t *Team) checkDuplicateTeamPlayers (field interface{}) error{
	players, ok := field.([]string)

	if !ok{
		return fmt.Errorf("could not parse %T into object", field)
	}
	
	//Create a map of string to int (value is irrelevant, I only care about the key)
	uniqueNames := make(map[string]struct{})

	//In order to see how if there are duplicate names, add all of the names the client passed into the map I just made.
	//Let's say players len = 5. If all 5 names are unique, that map will also have 5 elements. If there is one
	//duplicate (m, n, l, k, l), the map will have a length of 4 because the second "l" in that example is not given 
	//its own bucket in the map, the value of the first "l" is simply overridden.
	for _, player := range players{
		
		//turn all player names to lower, and remove all spaces to prevent anamolies like "alice" and " A L ICE "
		//from passing as non-duplicates.
		player = strings.ReplaceAll(strings.ToLower(player), " ", "")

		//Check if the editted name exists in the map. If it does, a duplicate was found.
		if _, exists := uniqueNames[player]; exists{
			return fmt.Errorf("no duplicate name allowed! %s", player)
		}

		//Afterwards, add the name to the map. On the second name, it will be compared to the first, the
		//third will be compared to the first and second, and so on.
		uniqueNames[player] = struct{}{}
	}

	//If not, return nil signifying validation success.
	return nil
}