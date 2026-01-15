package controllers

import (
	"JeopardyScoreBoardV2/services"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SavedGamesController struct {
	Router          *chi.Mux
	savedGameService services.SaveGameService
}

func (s *SavedGamesController) Init(service services.SaveGameService) {
	s.Router = chi.NewRouter()
	s.savedGameService = service

	s.Router.Get("/", s.GetAllSavedGames)
	s.Router.Post("/", s.AddSavedGame)
	s.Router.Get("/{location_name}", s.GetAllSavedGamesFromLocation)
	s.Router.Delete("/{location_name}", s.DeleteSavedGame)
}

func (s *SavedGamesController) GetAllSavedGames(res http.ResponseWriter, req *http.Request){

}

func (s *SavedGamesController) GetAllSavedGamesFromLocation(res http.ResponseWriter, req *http.Request){

}

func (s *SavedGamesController) AddSavedGame(res http.ResponseWriter, req *http.Request){

}

func (s *SavedGamesController) DeleteSavedGame(res http.ResponseWriter, req *http.Request){

}