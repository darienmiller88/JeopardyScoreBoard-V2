package services

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
	"JeopardyScoreBoardV2/utils"
	"net/http"
)

type SaveGameService interface{
	GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame]
	AddSavedGame(savedGame models.SavedGame)          models.Result[models.SavedGame]
	DeleteSavedGame(savedGameId string)               models.Result[string]
	GetAllSavedGames()                                models.Result[[]models.SavedGame]
	

}

type SaveGameServiceImpl struct{
	SavedGameRepository repositories.SavedGameRepository 
	LocationRepository  repositories.LocationRepository
	PlayerRepository    repositories.PlayerRepository
}

func (s *SaveGameServiceImpl) GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGamesFromLocationDB(locationName)
}

func (s *SaveGameServiceImpl) GetAllSavedGames() models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGamesDB()
}

//Business rules for a saved game:

/*
- Any player from any site can play at any site
- Winning player MUST exist if a player game is played
- All players who play MUST exist in the players Table
- game type can be either a team game or player game, but not both and not neither (validated in model)
- in a team game or player game, at least one team or player MUST be present
- a player game CANNOT have teams added
- a team game CANNOT have players added
- location id must exist in the locations table
- winning team id MUST exist if a team game is played
- total score must equal to the sum score of all players/teams
- average score must equal to the average score of all players/teams
*/
func (s *SaveGameServiceImpl) AddSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame]{
	//Validate the game to ensure it's either a team game or saved game, and that both have at
	//least one player or team participating.
	if err := savedGame.Validate(); err != nil{
		return utils.GetResult(err, http.StatusUnprocessableEntity, savedGame)
	}

	if savedGame.WinningPlayerId.Valid {
		//Check if the location id for the saved game exists
		if result := s.LocationRepository.IsLocationIdValid(savedGame.LocationId); result.Err != nil {
			return utils.GetResult(result.Err, result.StatusCode, savedGame)
		}

		//Check if the players the client added actually exist.
		if result := s.PlayerRepository.ArePlayersValid(savedGame.Players); result.Err != nil {
			return utils.GetResult(result.Err, result.StatusCode, savedGame)
		}
		
		//Check if the winning player exists
		if result := s.SavedGameRepository.IsWinningPlayerValid(savedGame.WinningPlayerName.String); result.Err != nil {
			return utils.GetResult(result.Err, result.StatusCode, savedGame)
		}


		savedGame.TotalPoints = s.getTotalScore(savedGame)
		savedGame.AveragePoints = s.getAverageScore(savedGame)
	} else{

	}



	//Check if the winning team exists
	if result := s.SavedGameRepository.IsWinningTeamIdValid(int(savedGame.WinningTeamId.Int32)); result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, savedGame)
	}

	return s.SavedGameRepository.AddSavedGameDB(savedGame)
}

func (s *SaveGameServiceImpl) DeleteSavedGame(savedGameId string) models.Result[string]{
	return s.SavedGameRepository.DeleteSavedGameDB(savedGameId)
}

func (s *SaveGameServiceImpl) getTotalScore(savedGame models.SavedGame) int{
	sum := 0
	
	if savedGame.WinningPlayerId.Valid {
		for _, player := range savedGame.Players {
			sum += player.Score
		}
	} else{
		for _, team := range savedGame.Teams {
			sum += team.Score
		}
	}

	return sum
}

func (s *SaveGameServiceImpl) getAverageScore(savedGame models.SavedGame) float64{
	averageScore := 0.0
	
	if savedGame.WinningPlayerId.Valid {
		averageScore = float64(s.getTotalScore(savedGame)) / float64(len(savedGame.Players))
	} else{
		averageScore = float64(s.getTotalScore(savedGame)) / float64(len(savedGame.Teams))		
	}

	return averageScore
}