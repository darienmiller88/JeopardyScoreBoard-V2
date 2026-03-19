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
	Teams        []models.Team
	Router       *chi.Mux
	template     *template.Template
	TeamService  services.TeamService
}

func (t *TeamsController) Init(teamService services.TeamService) {
	t.Router = chi.NewRouter()	

	//Add chi routes here
	t.Teams = append(t.Teams, models.Team{
		ID: 0,
		TeamName: "Flushing",
		Score: 0,
		PlayerNames: []string{
			"Michael Rosado",
			"Christopher Taylor",
			"Queen Branch",
			"Earnest Walker",
			"Jabriel Graham",
		},
	}, models.Team{
		ID: 1,
		TeamName: "Port Richmond",
		Score: 0,
		PlayerNames: []string{
			"Michael Melendez",
			"Margaret Dockery",
			"Cathleen Lee",
			"Betty Mitchell",
			"Schevon Williams",
			"Jimmy Ramirez",
			"Nijmah Othman",
			"Denise Farmer",
			"Michelle Tesoriero",
		},
	})

	t.Router.Get("/", t.GetTeams)
	t.Router.Get("/{id}", t.GetTeamPlayersByTeamId)
	t.Router.Post("/{id}/add-points", t.AddPoints)
	t.Router.Post("/{id}/minus-points", t.MinusPoints)


	templ, err := template.ParseGlob("templates/partials/*.html")

	if err != nil {
		panic(err)
	}

	t.template = templ
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

func (t *TeamsController) GetTeams(res http.ResponseWriter, req *http.Request) {
	if err := t.template.ExecuteTemplate(res, "team_cards", t.Teams); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TeamsController) GetTeamPlayersByTeamId(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	idInt, _ := strconv.Atoi(id)

	team := t.Teams[idInt]

	if err := t.template.ExecuteTemplate(res, "team_players", team.PlayerNames); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}