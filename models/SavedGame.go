package models

import (
	"fmt"
	"math"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type SavedGame struct {
	ID            int       `db:"id"` // The MongoDB document ID
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	LocationName  string    `db:"location_name"`
	Players       *[]Player      `bson:",omitempty"`
	Teams         *[]Team            `bson:",omitempty"`
	TotalPoints   int 			     `bson:"total_points"`
	AveragePoints float64            `bson:"average_points"`
	Winner        string    `db:"winner"`
}

// type SavedGame struct {
// 	ID            bson.ObjectID `bson:"_id,omitempty"` // The MongoDB document ID
// 	CreatedAt     time.Time          `bson:"created_at"`
// 	UpdatedAt     time.Time          `bson:"updated_at"`
// 	LocationName  string			 `bson:"location_name"`
// 	Players       *[]Player      `bson:",omitempty"`
// 	Teams         *[]Team            `bson:",omitempty"`
// 	TotalPoints   int 			     `bson:"total_points"`
// 	AveragePoints float64            `bson:"average_points"`
// 	Winner        Player         `bson:"winner"`
// }

func (s *SavedGame) Validate() error{
	return validation.ValidateStruct(
		s,

		//Order matters here! Validate the location name first to ensure that the Location field is not nil.
		validation.Field(&s.LocationName, validation.Required),

		//Validate the teams field if the user chooses to add it.
		validation.Field(
			&s.Teams, 

			//Validate to check to make sure that the client did not send in both a non-nil players AND teams field
			validation.By(s.validatePlayersAndTeams), 

			//Afterwards, enforce a "Required" requirement for the teams field when there is no Players field
			validation.Required.When(s.Players == nil), 

			//When both checks pass, validate each team the client sent 
			validation.By(s.validateTeams),
		),
		
		//Validate the players field if the user chooses to add it.
		validation.Field(
			&s.Players, 

			//Validate to check to make sure that the client did not send in both a non-nil players AND teams field
			validation.By(s.validatePlayersAndTeams), 

			//Afterwards, enforce a requirement for the players field when there is no Teams field
			validation.Required.When(s.Teams == nil).Error("Must include at least one player"), 
		),
	)
}

func (s *SavedGame) InitCreatedAtAndUpdatedAt(){
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *SavedGame) CalcAveragePoints(){
	//If the user is playing a team game, calculate the average points using the number of teams, otherwise calculate it
	//using the total number of people at each ADAPT location.
	if s.Teams != nil {
		s.AveragePoints = math.Round(float64(s.TotalPoints) / float64(len(*s.Teams)))
	}else{
		s.AveragePoints = math.Round(float64(s.TotalPoints) / float64(len(*s.Players)))
	}
}

//Ensure that players and teams aren't both included
func (s *SavedGame) validatePlayersAndTeams(field interface{}) error{
	if s.Players != nil && s.Teams != nil {
		return fmt.Errorf("field 'players' and field 'teams' both cannot be included")
	}

	return nil
}

func (s *SavedGame) validateTeams(field interface{}) error{
	if s.Teams != nil {
		teamLimit := 2
		
		//If the client includes the Teams field, ensure they include exactly 2.
		if len(*s.Teams) != teamLimit {
			return fmt.Errorf("please include exactly %d teams", teamLimit)
		} 

		uniqueTeams := make(map[string]int)

		//In order to see if there are duplciate teams, create a map out of the teams the client sent.
		for _, team := range *s.Teams{
			uniqueTeams[team.TeamName] = 0
		}

		//If the number of unique teams is less than the total number of teams, there are duplicates.
		if len(uniqueTeams) < len(*s.Teams) {
			return fmt.Errorf("no duplicate teams allowed")
		}

		//if they include the Teams field, and it has only 2 unique teams, validate each team to ensure each team is valid,
		//which entails a valid ADAPT location, and actual people there.
		for _, team := range *s.Teams {
			if err := team.Validate(); err != nil {
				return err
			}
		}
	} 

	return nil
}