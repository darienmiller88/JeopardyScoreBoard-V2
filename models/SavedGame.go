package models

import (
	"database/sql"
	"errors"
	"time"
)

type SavedGame struct {
	ID                int            `db:"id"` 
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
	TotalPoints       int 		     `db:"total_score"`
	AveragePoints     float64        `db:"average_score"`
	WinningPlayerName sql.NullString `db:"winning_player_name"`
	WinningTeamId     sql.NullInt32  `db:"winning_team_id"`
	WinningPlayerId   sql.NullInt32  `db:"winning_player_id"`
	LocationId        int            `json:"location_id" db:"location_id"`

	//This will be used to determine which array gets filled: Teams or Players
	IsPlayerGame      bool		     `json:"is_player_game"`
	
	//These fields will be used to fill WinningPlayerId, WinningTeamId, and WinningPlayerName. The less
	//I depend on the client for correct information, the better.
	Players           []Player       `json:"players" db:"-"`
	Teams             []Team         `json:"teams"   db:"-"`
}

func (s *SavedGame) Validate() error{
	if err := s.validateParticipation(); err != nil{
		return err
	}

	return nil
}

func (s *SavedGame) validateParticipation() error{
	//When a player game is being played, there can be no teams added
	if s.IsPlayerGame && len(s.Teams) != 0{
		return errors.New("a player game cannot have any teams added")
	}

	//When a team game is being played, there can be no players added
	if !s.IsPlayerGame && len(s.Players) != 0{
		return errors.New("a team game cannot have any players added")
	}

	//If a player game is being played, players MUST be filled with at least one player
	if s.IsPlayerGame && len(s.Players) == 0{
		return errors.New("players cannot be empty when winning player id is supplied")
	}

	//If a team game is being played, teams MUST be filled with at least one team
	if !s.IsPlayerGame && len(s.Teams) == 0{
		return errors.New("teams cannot be empty when winning team id is supplied")
	}

	return nil
}

//Calculate total score for both player games and team games
func (s *SavedGame) CalculateTotalPoints(){	
	s.TotalPoints = 0

	if s.IsPlayerGame {
		for _, player := range s.Players {
			s.TotalPoints += player.Score
		}
	} else{
		for _, team := range s.Teams {
			s.TotalPoints += team.Score
		}
	}
}

//Calculate average score for both player games and team games
func (s *SavedGame) CalculateAveragePoints(){
	s.CalculateTotalPoints()

	if s.IsPlayerGame {
		s.AveragePoints = float64(s.TotalPoints) / float64(len(s.Players))
	} else{
		s.AveragePoints = float64(s.TotalPoints) / float64(len(s.Teams))		
	}
}

//Calculate the winning team or player. If there is only one team or player, just set them as the 
//winning team or player. Otherwise, loop through and find the winner manually.
//At this point, it is assumed that the .Validate() method has been called first.
func (s *SavedGame) CalculateWinner(){
	if len(s.Players) == 1 {
		s.WinningPlayerName.String = s.Players[0].PlayerName
		s.WinningPlayerId.Int32 = int32(s.Players[0].ID)
	} else if len(s.Teams) == 1{
		s.WinningTeamId.Int32 = int32(s.Teams[0].ID)
	} else{
		if s.IsPlayerGame {
			s.calcWinnerForPlayers()
		} else{
			s.calcWinnerForTeams()
		}
	}
}

func (s *SavedGame) calcWinnerForPlayers(){
	winningPlayer := s.Players[0]
	
	for _, player := range s.Players[1:] {
		if player.Score > winningPlayer.Score {
			winningPlayer = player
		}
	}

	s.WinningPlayerName.String = winningPlayer.PlayerName
	s.WinningPlayerId.Int32 = int32(winningPlayer.ID)
}

func (s *SavedGame) calcWinnerForTeams(){
	winningTeam := s.Teams[0]
	
	for _, team := range s.Teams[1:] {
		if team.Score > winningTeam.Score {
			winningTeam = team
		}
	}	
	
	s.WinningTeamId.Int32 = int32(winningTeam.ID)
}