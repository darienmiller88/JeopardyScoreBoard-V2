package controllers

import (
	// "fmt"
	// "errors"
	// "net/http"
	// "net/http/httptest"
	// "testing"

	// "github.com/go-chi/chi/v5"
	// "github.com/stretchr/testify/assert"

	"JeopardyScoreBoardV2/models"
)

type mockPlayerService struct{
	playersResult models.Result[[]models.Player]
	playerResult  models.Result[models.Player]
}

func(m *mockPlayerService) UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]{
	return m.playerResult
}

func (m *mockPlayerService) AddPlayerToLocation(locationName string, playerName string)  models.Result[models.Player] {
	return m.playerResult
}	
	
func (m *mockPlayerService) RemovePlayer(playerName string) models.Result[models.Player] {
	return m.playerResult
}
	
func (m *mockPlayerService) GetPlayersFromLocation(locationName string) models.Result[[]models.Player] {
	return m.playersResult
}
	 
func (m *mockPlayerService) GetAllPlayersFromAllLocations() models.Result[[]models.Player]{
	return m.playersResult
}	