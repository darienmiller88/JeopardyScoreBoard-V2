package controllers

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ViewsController struct{
	Router           *chi.Mux
	pagesTemplate    *template.Template
	partialsTemplate *template.Template
}

func (v *ViewsController) Init(){
	v.Router = chi.NewRouter()
	v.pagesTemplate  = template.Must(template.ParseGlob("templates/*.html"))
	v.partialsTemplate = template.Must(template.ParseGlob("templates/partials/*.html"))

	v.Router.Get("/", v.CreateGame)
	v.Router.Get("/team-mode", v.TeamMode)
}

func (v *ViewsController) CreateGame(res http.ResponseWriter, req *http.Request){
	if err := v.pagesTemplate.ExecuteTemplate(res, "CreateGame.html", nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) TeamMode(res http.ResponseWriter, req *http.Request){
	if err := v.pagesTemplate.ExecuteTemplate(res, "TeamMode.html", nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) AddPlayer(res http.ResponseWriter, req *http.Request){
	if err := v.pagesTemplate.ExecuteTemplate(res, "AddPlayer.html", nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}