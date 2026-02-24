package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"

	"JeopardyScoreBoardV2/services"
)

// Will serve as HTMX end points to Add players, Delete players and Update their names
type PlayersController struct {
	Router        *chi.Mux
	template      *template.Template
	playerService services.PlayerService
}

func (p *PlayersController) Init(service services.PlayerService) {
	p.Router = chi.NewRouter()
	p.playerService = service

	p.Router.Get("/by-location", p.GetAllPlayersFromOneLocation)
	p.Router.Put("/{location_name}", p.UpdatePlayerName)
	p.Router.Post("/", p.AddPlayerToLocation)
	p.Router.Delete("/{location_name}", p.RemovePlayer)

	t, err := template.ParseGlob("templates/partials/*.html")

	if err != nil {
		panic(err)
	}

	p.template = t
}

func (p *PlayersController) GetAllPlayersFromOneLocation(res http.ResponseWriter, req *http.Request) {
	location := req.URL.Query().Get("location")
	result := p.playerService.GetPlayersFromLocation(location)

	fmt.Println("location:", location)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	data := template.FuncMap{
		"Players": result.ResultData,
		"SelectedLocation": location,
		"HasPlayers": len(result.ResultData) > 0,
	}

	if err := p.template.ExecuteTemplate(res, "player_list_section", data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (p *PlayersController) AddPlayerToLocation(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	locationName := req.FormValue("location")
	firstName := req.FormValue("first")
	lastName := req.FormValue("last")
	result := p.playerService.AddPlayerToLocation(locationName, firstName, lastName)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	if err := p.template.ExecuteTemplate(res, "player_card", result.ResultData); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (p *PlayersController) RemovePlayer(res http.ResponseWriter, req *http.Request) {
	locationName := chi.URLParam(req, "location_name")
	playerName := ""

	if err := json.NewDecoder(req.Body).Decode(&playerName); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	result := p.playerService.RemovePlayer(playerName, locationName)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(result)
}

func (p *PlayersController) UpdatePlayerName(res http.ResponseWriter, req *http.Request) {
	locationName := chi.URLParam(req, "location_name")
	names := struct {
		OldPlayerName string `json:"old_player_name"`
		NewPlayerName string `json:"new_player_name"`
	}{}

	if err := json.NewDecoder(req.Body).Decode(&names); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	result := p.playerService.UpdatePlayerName(names.OldPlayerName, names.NewPlayerName, "", locationName)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(result)
}
