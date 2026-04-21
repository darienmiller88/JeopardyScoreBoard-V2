package controllers

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/services"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TeamsController struct {
	Teams       []models.Team
	Router      *chi.Mux
	template    *template.Template
	TeamService services.TeamService
}

func (t *TeamsController) Init(teamService services.TeamService) {
	t.Router = chi.NewRouter()

	t.Router.Get("/", t.GetTeams)
	t.Router.Get("/team-names", t.GetTeamNames)
	t.Router.Get("/{id}", t.GetTeamPlayersByTeamId)
	t.Router.Post("/{id}/add-points", t.AddPoints)
	t.Router.Post("/{id}/minus-points", t.MinusPoints)
	t.Router.Post("/add-team", t.AddTeam)
	t.Router.Delete("/{id}", t.DeleteTeam)

	templ, err := template.ParseGlob("templates/partials/*.html")

	if err != nil {
		panic(err)
	}

	t.template = templ
}

func (t *TeamsController) DeleteTeam(res http.ResponseWriter, req *http.Request){
	id := chi.URLParam(req, "id")
	idInt, _ := strconv.Atoi(id)
	t.Teams = append(t.Teams[:idInt], t.Teams[idInt+1:]...)

	res.WriteHeader(http.StatusOK)
}

func (t *TeamsController) AddTeam(res http.ResponseWriter, req *http.Request) {
	teamName := req.FormValue("teamName")

	if teamName == "" {
		http.Error(res, "Invalid team name", http.StatusBadRequest)
		return
	}

	team := models.Team{
		ID:       len(t.Teams),
		Score:    0,
		TeamName: teamName,
	}

	t.Teams = append(t.Teams, team)

	if err := t.template.ExecuteTemplate(res, "team_card", team); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *TeamsController) AddPoints(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	idInt, _ := strconv.Atoi(id)
	points, _ := strconv.Atoi(req.FormValue("points"))

	// fallback if input is empty
	if points == 0 {
		points = 100
	}

	t.Teams[idInt].Score += points

	t.template.ExecuteTemplate(res, "score", struct {
		ID    int
		Score int
	}{idInt, t.Teams[idInt].Score})
}

func (t *TeamsController) MinusPoints(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	idInt, _ := strconv.Atoi(id)

	points, _ := strconv.Atoi(req.FormValue("points"))
	if points == 0 {
		points = 100
	}

	t.Teams[idInt].Score -= points

	t.template.ExecuteTemplate(res, "score", struct {
		ID    int
		Score int
	}{idInt, t.Teams[idInt].Score})
}

func (t *TeamsController) GetTeamNames(res http.ResponseWriter, req *http.Request) {
	if err := t.template.ExecuteTemplate(res, "team_cards", t.Teams); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TeamsController) GetTeams(res http.ResponseWriter, req *http.Request) {
	if err := t.template.ExecuteTemplate(res, "team_cards", t.Teams); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TeamsController) GetTeamPlayersByTeamId(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	idInt, _ := strconv.Atoi(id)

	team := t.Teams[idInt]

	if err := t.template.ExecuteTemplate(res, "team_players", team.PlayerNames); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
