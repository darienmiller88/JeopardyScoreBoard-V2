package controllers

import (
	"JeopardyScoreBoardV2/services"
	"html/template"

	"github.com/go-chi/chi/v5"
)

type TeamsController struct {
	Router           *chi.Mux
	template         *template.Template
	TeamService      services.TeamService
}

func (t *TeamsController) Init(teamService services.TeamService) {
	t.Router = chi.NewRouter()	
	templ, err := template.ParseGlob("templates/partials/*.html")

	if err != nil {
		panic(err)
	}

	t.template = templ
}