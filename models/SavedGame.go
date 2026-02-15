package models

import (
	"JeopardyScoreBoardV2/encryption"
	"database/sql"
	"errors"
	"time"
)

type SavedGame struct {
	ID                         int           `db:"id"` 
	CreatedAt                  time.Time     `db:"created_at"`
	UpdatedAt                  time.Time     `db:"updated_at"`
	TotalPoints       		   int 		     `db:"total_score"`
	AveragePoints    		   float64       `db:"average_score"`
	WinningPlayerNameEncrypted []byte        `db:"winning_player_name_encrypted"`
	WinningPlayerNameHash      []byte        `db:"winning_player_name_hash"`
	WinningTeamId              sql.NullInt32 `db:"winning_team_id"`
	WinningPlayerId            sql.NullInt32 `db:"winning_player_id"`
	LocationId                 int           `json:"location_id" db:"location_id"`

	//This will be used to determine which array gets filled: Teams or Players
	IsPlayerGame      bool	  `json:"is_player_game" db:"-"`
	
	//Will contained the decrypted winning player name
	WinningPlayerName string  `json:"winning_player_name"`
	
	//These fields will be used to fill WinningPlayerId, WinningTeamId, and WinningPlayerName. The less
	//I depend on the client for correct information, the better. Therefore, these will be the only
	//fields to be validated.
	Players           []Player `json:"players" db:"-"`
	Teams             []Team   `json:"teams"   db:"-"`

	//encryption service to encrypt winning player name
	encryptionService *encryption.EncryptionService
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
		return errors.New("players cannot be empty when a player game is being played")
	}

	//If a team game is being played, teams MUST be filled with at least one team
	if !s.IsPlayerGame && len(s.Teams) == 0{
		return errors.New("teams cannot be empty when a team game is being played")
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
func (s *SavedGame) CalculateWinner() error{
	if len(s.Players) == 1 {
		encryptedName, err := s.encryptionService.Encrypt(s.Players[0].PlayerName)

		if err != nil {
			return err
		}

		s.WinningPlayerNameEncrypted = encryptedName
		s.WinningPlayerId = newNullInt32(s.Players[0].ID)
	} else if len(s.Teams) == 1{
		s.WinningTeamId = newNullInt32(s.Teams[0].ID)
	} else{
		if s.IsPlayerGame {
			if err := s.calcWinnerForPlayers(); err != nil{
				return err
			}
		} else{
			s.calcWinnerForTeams()
		}
	}

	return nil
}

func (s *SavedGame) calcWinnerForPlayers() error{
	winningPlayer := s.Players[0]
	
	for _, player := range s.Players[1:] {
		if player.Score > winningPlayer.Score {
			winningPlayer = player
		}
	}

	playerNameDecrypted, err := s.encryptionService.Decrypt(winningPlayer.PlayerNameEncrypted)

	if err != nil {
		return err
	}

	//Encrypt and hash the winning player name
	s.WinningPlayerNameEncrypted = winningPlayer.PlayerNameEncrypted
	s.WinningPlayerNameHash =  encryption.NameHash(winningPlayer.PlayerName)
	s.WinningPlayerName = playerNameDecrypted
	s.WinningPlayerId = newNullInt32(winningPlayer.ID)

	return nil
}

func (s *SavedGame) calcWinnerForTeams(){
	winningTeam := s.Teams[0]
	
	for _, team := range s.Teams[1:] {
		if team.Score > winningTeam.Score {
			winningTeam = team
		}
	}	
	
	s.WinningTeamId = newNullInt32(winningTeam.ID)
}

func newNullInt32(newInt int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(newInt), Valid: true}
}

func newNullString(newString string) sql.NullString {
	return sql.NullString{String: newString, Valid: true}
}
