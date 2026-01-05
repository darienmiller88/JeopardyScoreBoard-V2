package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"JeopardyScoreBoardV2/services"
)

type PlayersController struct{
	Router *chi.Mux
	playerService services.PlayerService
}

func (p *PlayersController) Init(service services.PlayerService){
	p.Router = chi.NewRouter()
	p.playerService = service

	p.Router.Get("/", p.GetAllPlayers)
	p.Router.Get("/{location_name}", p.GetAllPlayersFromOneLocation)
	p.Router.Post("/{location_name}", p.AddPlayerToLocation)
	p.Router.Put("/", p.UpdatePlayerName)
}

func (p *PlayersController) GetAllPlayers(res http.ResponseWriter, req *http.Request){
	result := p.playerService.GetAllPlayersFromAllLocations()

	if result.Err != nil{
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(result)
}

func (p *PlayersController) GetAllPlayersFromOneLocation(res http.ResponseWriter, req *http.Request){

}

func (p *PlayersController) AddPlayerToLocation(res http.ResponseWriter, req *http.Request){
	locationName := chi.URLParam(req, "location_name")
	playerName := ""
	
	if err := json.NewDecoder(req.Body).Decode(&playerName); err !=nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	result := p.playerService.AddPlayerToLocation(locationName, playerName)

	if result.Err != nil{
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	data, err := json.Marshal(&result)

	if err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	
	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	res.Write(data)
}

func (p *PlayersController) RemovePlayerFromLocation(res http.ResponseWriter, req *http.Request){

}

func (p *PlayersController) UpdatePlayerName(res http.ResponseWriter, req *http.Request){
	names := struct{
		OldPlayerName string `json:"old_player_name"`
		NewPlayerName string `json:"new_player_name"`
	}{}
	
	if err := json.NewDecoder(req.Body).Decode(&names); err !=nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	result := p.playerService.UpdatePlayerName(names.OldPlayerName, names.NewPlayerName)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}
	
	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(result)
}
