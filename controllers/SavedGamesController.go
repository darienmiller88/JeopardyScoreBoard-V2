package controllers

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SavedGamesController struct {
	Router          *chi.Mux
	savedGameService services.SaveGameService
}

func (s *SavedGamesController) Init(service services.SaveGameService) {
	s.Router = chi.NewRouter()
	s.savedGameService = service

	s.Router.Get("/{id}/players", s.GetAllPlayersFromSavedGame)
	s.Router.Get("/", s.GetAllSavedGames)
	s.Router.Post("/", s.AddSavedGame)
	s.Router.Get("/{location_name}", s.GetAllSavedGamesFromLocation)
	s.Router.Delete("/{id}", s.DeleteSavedGame)
}

func (s *SavedGamesController) GetAllPlayersFromSavedGame(res http.ResponseWriter, req *http.Request){
	
}

func (s *SavedGamesController) GetAllSavedGames(res http.ResponseWriter, req *http.Request){
	result := s.savedGameService.GetAllSavedGames()

	if result.Err != nil{
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(result)
}

func (s *SavedGamesController) GetAllSavedGamesFromLocation(res http.ResponseWriter, req *http.Request){
	locationName := chi.URLParam(req, "location_name")
	result := s.savedGameService.GetAllSavedGamesFromLocation(locationName)

	if result.Err != nil{
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(result)
}

func (s *SavedGamesController) AddSavedGame(res http.ResponseWriter, req *http.Request){
	savedGame := models.SavedGame{}

	if err := json.NewDecoder(req.Body).Decode(&savedGame); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	result := s.savedGameService.AddSavedGame(savedGame)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(result)
}

func (s *SavedGamesController) DeleteSavedGame(res http.ResponseWriter, req *http.Request){
	id := chi.URLParam(req, "id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	result := s.savedGameService.DeleteSavedGame(idInt)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(result)
}