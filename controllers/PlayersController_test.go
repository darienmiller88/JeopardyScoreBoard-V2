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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlayerService struct {
	playersResult models.Result[[]models.Player]
	playerResult  models.Result[models.Player]
}

func (m *mockPlayerService) UpdatePlayerName(oldPlayerName string, newPlayerName string, locationName string) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerService) AddPlayerToLocation(locationName string, playerName string) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerService) RemovePlayer(playerName string, locationName string) models.Result[models.Player] {
	return m.playerResult
}

func (m *mockPlayerService) GetPlayersFromLocation(locationName string) models.Result[[]models.Player] {
	return m.playersResult
}

func (m *mockPlayerService) GetAllPlayersFromAllLocations() models.Result[[]models.Player] {
	return m.playersResult
}

func getPlayersController(service services.PlayerService) *chi.Mux {
	playersController := PlayersController{}

	playersController.Init(service)

	return playersController.Router
}

////////////////////////////////
//CREATE/POST tests
///////////////////////////////

func TestAddPlayer_Ok(t *testing.T) {
	name := "Jane Doe"
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			StatusCode: http.StatusOK,
			ResultData: models.Player{PlayerName: name },
		},
	})

	body, _ := json.Marshal(name)
	_, rec := getResponseFromController(router, "/Elmwood", http.MethodPost, bytes.NewBuffer(body))

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var result models.Result[models.Player]
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, name, result.ResultData.PlayerName)
}

func TestAddPlayer_InvalidJSON(t *testing.T) {
	router := getPlayersController(&mockPlayerService{})

	_, rec := getResponseFromController(router, "/Elmwood", http.MethodPost, bytes.NewBufferString("{bad json"))

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid")
}

func TestAddPlayer_ValidationError(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("Player name must have exactly two parts"),
			StatusCode: http.StatusUnprocessableEntity,
		},
	})

	body, _ := json.Marshal("Jane")

	_, rec := getResponseFromController(router, "/Elmwood", http.MethodPost, bytes.NewBuffer(body))

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "two parts")
}

func TestAddPlayer_LocationNotFound(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("no location 'Nowhere' found"),
			StatusCode: http.StatusNotFound,
		},
	})

	body, _ := json.Marshal("Jane Doe")
	_, rec := getResponseFromController(router, "/Nowhere", http.MethodPost, bytes.NewBuffer(body))

	require.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "no location")
}

func TestAddPlayer_NameTaken(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("name Jane Doe is already taken"),
			StatusCode: http.StatusConflict,
		},
	})

	body, _ := json.Marshal("Jane Doe")

	_, rec := getResponseFromController(router, "/Elmwood", http.MethodPost, bytes.NewBuffer(body))

	require.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "already taken")
}

func TestAddPlayer_DbError(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("db down"),
			StatusCode: http.StatusInternalServerError,
		},
	})

	body, _ := json.Marshal("Jane Doe")
	_, rec := getResponseFromController(router, "/Elmwood", http.MethodPost, bytes.NewBuffer(body))

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "db down")
}






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

	resp, rec := getResponseFromController(router, "/", http.MethodGet, nil)

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

	resp, rec := getResponseFromController(router, "/", http.MethodGet, nil)

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

	resp, rec := getResponseFromController(router, "/", http.MethodGet, nil)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.Result[[]models.Player]
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Empty(t, result.ResultData)
}

func TestGetPlayersFromLocation_Ok(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playersResult: models.Result[[]models.Player]{
			StatusCode: http.StatusOK,
			ResultData: []models.Player{
				{PlayerName: "Jane Doe"},
				{PlayerName: "John Smith"},
			},
		},
	})

	_, rec := getResponseFromController(router, "/Elmwood", http.MethodGet, nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var result models.Result[[]models.Player]
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Len(t, result.ResultData, 2)
}

func TestGetPlayersFromLocation_Empty(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playersResult: models.Result[[]models.Player]{
			StatusCode: http.StatusOK,
			ResultData: []models.Player{},
		},
	})

	_, rec := getResponseFromController(router, "/Elmwood", http.MethodGet, nil)

	require.Equal(t, http.StatusOK, rec.Code)

	var result models.Result[[]models.Player]
	err := json.Unmarshal(rec.Body.Bytes(), &result)

	//Ensure the controller returned no error, and the empty list of players
	require.NoError(t, err)
	assert.Empty(t, result.ResultData)
}

func TestGetPlayersFromLocation_DBError(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playersResult: models.Result[[]models.Player]{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("db error"),
		},
	})

	resp, rec := getResponseFromController(router, "/Elmwood", http.MethodGet, nil)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, rec.Body.String(), "db error")
}








////////////////////////////////
//UPDATE/PUT tests
///////////////////////////////

func TestUpdatePlayerName_Ok(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			StatusCode: http.StatusOK,
			ResultData: models.Player{PlayerName: "Jane Doe"},
		},
	})

	body := `{"old_player_name":"John Doe","new_player_name":"Jane Doe"}`
	_, rec := getResponseFromController(router, "/", http.MethodPut, bytes.NewBufferString(body))

	require.Equal(t, http.StatusOK, rec.Code)

	var result models.Result[models.Player]
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, "Jane Doe", result.ResultData.PlayerName)
}

func TestUpdatePlayerName_InvalidJSON(t *testing.T) {
	router := getPlayersController(&mockPlayerService{})
	_, rec := getResponseFromController(router, "/", http.MethodPut, bytes.NewBufferString("{bad json"))

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid")
}

func TestUpdatePlayerName_ValidationError(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("validation failed"),
			StatusCode: http.StatusUnprocessableEntity,
		},
	})

	body := `{"old_player_name":"","new_player_name":""}`
	_, rec := getResponseFromController(router, "/", http.MethodPut, bytes.NewBufferString(body))

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "validation")
}

func TestUpdatePlayerName_NameTaken(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("name already exists"),
			StatusCode: http.StatusConflict,
		},
	})

	body := `{"old_player_name":"John","new_player_name":"Jane"}`
	_, rec := getResponseFromController(router, "/", http.MethodPut, bytes.NewBufferString(body))

	require.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "name")
}

func TestUpdatePlayerName_OldNameNotFound(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("player not found"),
			StatusCode: http.StatusNotFound,
		},
	})

	body := `{"old_player_name":"Ghost","new_player_name":"Jane"}`
	_, rec := getResponseFromController(router, "/", http.MethodPut, bytes.NewBufferString(body))

	require.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "player")
}

func TestUpdatePlayerName_DbError(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("db down"),
			StatusCode: http.StatusInternalServerError,
		},
	})

	body := `{"old_player_name":"John","new_player_name":"Jane"}`
	_, rec := getResponseFromController(router, "/", http.MethodPut, bytes.NewBufferString(body))

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "db down")
}


////////////////////////////////
//DESTROY/DELETE tests
///////////////////////////////

func TestRemovePlayer_Ok(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			StatusCode: http.StatusOK,
			ResultData: models.Player{PlayerName: "Jane Doe"},
		},
	})

	_, rec := getResponseFromController(router, "/Jane%20Doe", http.MethodDelete, nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var result models.Result[models.Player]
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, "Jane Doe", result.ResultData.PlayerName)
}

func TestRemovePlayer_NotFound(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("player not found"),
			StatusCode: http.StatusNotFound,
		},
	})

	_, rec := getResponseFromController(router, "/FakePlayer", http.MethodDelete, nil)

	require.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "player not found")
}

func TestRemovePlayer_DbError(t *testing.T) {
	router := getPlayersController(&mockPlayerService{
		playerResult: models.Result[models.Player]{
			Err:        fmt.Errorf("db down"),
			StatusCode: http.StatusInternalServerError,
		},
	})

	_, rec := getResponseFromController(router, "/Jane%20Doe", http.MethodDelete, nil)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "db down")
}