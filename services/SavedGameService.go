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
	Repository repositories.SavedGameRepository
}

func (s *SaveGameServiceImpl) GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame]{
	return s.Repository.GetAllSavedGamesFromLocationDB(locationName)
}

func (s *SaveGameServiceImpl) GetAllSavedGames() models.Result[[]models.SavedGame]{
	return s.Repository.GetAllSavedGamesDB()
}

//Business rules for a saved game:

/*
- Any player from any site can play at any site
- Winning player MUST exist if a player game is played
- All players who play MUST exist in the players Table
- game type can be either a team game or player game, but not both and not neither (validated in model)
- in a team game or player game, at least one team or player MUST be present
- Team location id must exist in the locations table
- winning team id MUST exist if a team game is played
*/
func (s *SaveGameServiceImpl) AddSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame]{
	if err := savedGame.Validate(); err != nil{
		return utils.GetResult(err, http.StatusUnprocessableEntity, savedGame)
	}

	//Check if the players the client added actually exist.
	if result := s.Repository.ArePlayersValid(savedGame.Players); result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, savedGame)
	}

	//Check if the winning player exists
	if result := s.Repository.IsWinningPlayerValid(savedGame.WinningPlayerName.String); result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, savedGame)
	}

	//Check if the winning team exists
	if result := s.Repository.IsWinningTeamIdValid(int(savedGame.WinningTeamId.Int32)); result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, savedGame)
	}

	return s.Repository.AddSavedGameDB(savedGame)
}

func (s *SaveGameServiceImpl) DeleteSavedGame(savedGameId string) models.Result[string]{
	return s.Repository.DeleteSavedGameDB(savedGameId)
}