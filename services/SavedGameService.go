package services

import (
	"net/http"
	
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
)

type SaveGameService interface{

}

type SaveGameServiceImpl struct{
	SavedGameRepository repositories.SavedGameRepository
	LocationRepository  repositories.LocationRepository
}

func (s *SaveGameService) GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGamesFromLocation(locationName)
}

func (s *SaveGameService) GetAllSavedGames() models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGames()
}

func (s *SaveGameService) AddSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame]{
	return s.SavedGameRepository.AddSavedGame(savedGame)
}

func (s *SaveGameService) DeleteSavedGame(savedGameId int) models.Result[string]{
	return s.SavedGameRepository.DeleteSavedGame(savedGameId)
}