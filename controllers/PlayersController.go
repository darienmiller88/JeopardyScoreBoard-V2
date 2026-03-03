package controllers

import (
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

	p.Router.Get("/", p.GetAllPlayers)
	p.Router.Get("/by-location", p.GetAllPlayersFromOneLocation)
	p.Router.Put("/", p.UpdatePlayerName)
	p.Router.Delete("/", p.RemovePlayer)
	p.Router.Post("/", p.AddPlayerToLocation)

	t, err := template.ParseGlob("templates/partials/*.html")

	if err != nil {
		panic(err)
	}

	p.template = t
}

func (p *PlayersController) GetAllPlayers(res http.ResponseWriter, req *http.Request){
	result := p.playerService.GetAllPlayersFromAllLocations()

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	data := template.FuncMap{
		"Players":          result.ResultData,
		"HasPlayers":       len(result.ResultData) > 0,
		"SelectedLocation": "Elmwood",
	}

	p.template.ExecuteTemplate(res, "player_list_section", data)
}

func (p *PlayersController) GetAllPlayersFromOneLocation(res http.ResponseWriter, req *http.Request) {
	location := req.URL.Query().Get("location")
	result := p.playerService.GetPlayersFromLocation(location)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	data := template.FuncMap{
		"Players":          result.ResultData,
		"HasPlayers":       len(result.ResultData) > 0,
		"SelectedLocation": location,
	}

	if err := p.template.ExecuteTemplate(res, "player_list_section", data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (p *PlayersController) AddPlayerToLocation(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	locationName := req.FormValue("location")
	firstName := req.FormValue("first")
	lastName := req.FormValue("last")

	result := p.playerService.AddPlayerToLocation(locationName, firstName, lastName)
	data := template.FuncMap{
		"Player":           result.ResultData,
		"SelectedLocation": locationName,
	}

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	if err := p.template.ExecuteTemplate(res, "player_card", data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (p *PlayersController) RemovePlayer(res http.ResponseWriter, req *http.Request) {
	location := req.URL.Query().Get("location")
	playerId := req.URL.Query().Get("id")
	result := p.playerService.RemovePlayer(playerId, location)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	res.WriteHeader(200)
}

func (p *PlayersController) UpdatePlayerName(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	location := req.FormValue("location")
	playerId := req.FormValue("id")
	firstName := req.FormValue("first")
	lastName := req.FormValue("last")
	result := p.playerService.UpdatePlayerName(playerId, firstName, lastName, location)

	data := template.FuncMap{
		"Player":           result.ResultData,
		"SelectedLocation": location,
	}

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	if err := p.template.ExecuteTemplate(res, "player_card", data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
