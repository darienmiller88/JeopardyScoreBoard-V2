package services

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
)

type SaveGameService interface{
	GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame]
	AddSavedGame(savedGame models.SavedGame)          models.Result[models.SavedGame]
	DeleteSavedGame(savedGameId int)                  models.Result[string]
	GetAllSavedGames()                                models.Result[[]models.SavedGame]
}

type SaveGameServiceImpl struct{
	SavedGameRepository repositories.SavedGameRepository
	LocationRepository  repositories.LocationRepository
}

func (s *SaveGameServiceImpl) GetAllSavedGamesFromLocation(locationName string) models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGamesFromLocationDB(locationName)
}

func (s *SaveGameServiceImpl) GetAllSavedGames() models.Result[[]models.SavedGame]{
	return s.SavedGameRepository.GetAllSavedGamesDB()
}

func (s *SaveGameServiceImpl) AddSavedGame(savedGame models.SavedGame) models.Result[models.SavedGame]{
	return s.SavedGameRepository.AddSavedGameDB(savedGame)
}

func (s *SaveGameServiceImpl) DeleteSavedGame(savedGameId int) models.Result[string]{
	return s.SavedGameRepository.DeleteSavedGameDB(savedGameId)
}