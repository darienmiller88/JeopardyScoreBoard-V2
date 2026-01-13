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
	"JeopardyScoreBoardV2/services"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func getPlayersController(service services.PlayerService) *chi.Mux{
	playersController := PlayersController{}

	playersController.Init(service)

	return playersController.Router
}

////////////////////////////////
//CREATE/POST tests
///////////////////////////////


////////////////////////////////
//READ/GET tests
///////////////////////////////

func TestGetAllPlayers_Ok(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playersResult: models.Result[[]models.Player]{
			StatusCode: http.StatusOK,
			ResultData: []models.Player{
				{PlayerName: "Jane Doe"},
				{PlayerName: "John Smith"},
			},
		},
	})

	resp, rec := getResponseFromController(router, "/players")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result models.Result[[]models.Player]
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Len(t, result.ResultData, 2)
	assert.Equal(t, "Jane Doe", result.ResultData[0].PlayerName)
}

func TestGetAllPlayers_ServiceError(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playersResult: models.Result[[]models.Player]{
			Err:        fmt.Errorf("database exploded"),
			StatusCode: http.StatusInternalServerError,
		},
	})

	resp, rec := getResponseFromController(router, "/players")

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, rec.Body.String(), "database exploded")
}

func TestGetAllPlayers_EmptyList(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playersResult: models.Result[[]models.Player]{
			StatusCode: http.StatusOK,
			ResultData: []models.Player{},
		},
	})

	resp, rec := getResponseFromController(router, "/players")

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.Result[[]models.Player]
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Empty(t, result.ResultData)
}

////////////////////////////////
//UPDATE/PUT tests
///////////////////////////////




////////////////////////////////
//DESTROY/DELETE tests
///////////////////////////////


