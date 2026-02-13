package services

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
	"JeopardyScoreBoardV2/utils"
	"fmt"
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
	TeamRepository      repositories.TeamRepository
}

func (s *SaveGameServiceImpl) GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGamesFromLocationDB(locationName)
}

func (s *SaveGameServiceImpl) GetAllSavedGames() models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGamesDB()
}

//Business rules for a saved game:

/*
(set by saved game model)
- game type must be a player game or team game, but not both.
- player game must have at least one player 
- player game cannot have any teams added
- team game must have at least one team
- team game cannot have any players added

(set by saved game service)
- location id must exist in the locations table
- players must actually exist
- teams must actually exist (id must be real)

Winners for each type of game, and the total score will be calculated server side.
*/
func (s *SaveGameServiceImpl) AddSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame]{
	//Validate the game to ensure it's either a team game or saved game, and that both have at
	//least one player or team participating.
	if err := savedGame.Validate(); err != nil{
		return utils.GetResult(err, http.StatusUnprocessableEntity, savedGame)
	}

	//Check if the location id for the saved game exists
	if result := s.isLocationIdValid(savedGame.LocationId); result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, savedGame)
	}
	
	if savedGame.IsPlayerGame {
		//Check if the players the client added actually exist.
		result := s.arePlayersValid(savedGame.Players)
		
		if result.Err != nil {
			return utils.GetResult(result.Err, result.StatusCode, savedGame)
		}

		savedGame.Players = result.ResultData
	} else{
		//check if teams the client added actually exist
		result := s.areTeamsValid(savedGame.Teams)

		if result.Err != nil {
			return utils.GetResult(result.Err, result.StatusCode, savedGame)
		}

		savedGame.Teams = result.ResultData
	}
		
	//Perform the following calculations on the saved game
	savedGame.CalculateTotalPoints()
	savedGame.CalculateAveragePoints()
	savedGame.CalculateWinner()

	//Finaly, after validating for all edge cases, pass the saved game to repository to be inserted.
	return s.SavedGameRepository.AddSavedGameDB(savedGame)
}

func (s *SaveGameServiceImpl) DeleteSavedGame(savedGameId string) models.Result[string]{
	//check to see if saved game id exists first

	return s.SavedGameRepository.DeleteSavedGameDB(savedGameId)
}

//Checks if an location id actually exists
func (s *SaveGameServiceImpl) isLocationIdValid(locationId int) models.Result[models.Location]{
	result := s.LocationRepository.GetLocationById(locationId)

	if result.Err != nil {
		return result
	}

	return utils.GetResult(nil, http.StatusOK, result.ResultData)
}


//checks if a list of players is valid (exists in the database)
func (s *SaveGameServiceImpl) arePlayersValid(players []models.Player) models.Result[[]models.Player] {
    names := make([]string, len(players))

    for i, player := range players {
        names[i] = player.PlayerName
    }

    result := s.PlayerRepository.GetPlayersByNames(names)

    if result.Err != nil {
        return result
    }

    if len(result.ResultData) != len(names) {
        return utils.GetResult(fmt.Errorf("one or more players do not exist"), http.StatusNotFound, []models.Player{})
    }

    return result 
}

//Checks if teams are valid by seeing if their ids exist.
func (s *SaveGameServiceImpl) areTeamsValid(teams []models.Team) models.Result[[]models.Team]{
	teamIds := make([]int, len(teams))

    for i, team := range teams {
        teamIds[i] = team.ID
    }

    result := s.TeamRepository.GetAllTeamsByIds(teamIds)

    if result.Err != nil {
        return result
    }

    if len(result.ResultData) != len(teamIds) {
        return utils.GetResult(fmt.Errorf("one or more teams do not exist"), http.StatusNotFound, []models.Team{})
    }

    return result 
}