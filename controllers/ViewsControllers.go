package controllers

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ViewsController struct{
	Router           *chi.Mux
	pagesTemplate    *template.Template
}

func (v *ViewsController) Init(){
	v.Router = chi.NewRouter()
	v.pagesTemplate  = template.Must(template.ParseGlob("templates/*.html"))

	v.Router.Get("/", v.CreateGame)
	v.Router.Get("/team-mode", v.TeamMode)
	v.Router.Get("/add-player", v.AddPlayer)
	v.Router.Get("/view-games", v.ViewGames)
	v.Router.Get("/log-in", v.LogIn)
	v.Router.NotFound(v.NotFound)
}

func (v *ViewsController) CreateGame(res http.ResponseWriter, req *http.Request){
	// t, _ := template.ParseFiles("templates/Base.html", "templates/CreateGame.html")
	
	if err := v.pagesTemplate.ExecuteTemplate(res, "Base.html", nil); err != nil{
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

func (v *ViewsController) ViewGames(res http.ResponseWriter, req *http.Request){
	if err := v.pagesTemplate.ExecuteTemplate(res, "ViewGames.html", nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) LogIn(res http.ResponseWriter, req *http.Request){
	if err := v.pagesTemplate.ExecuteTemplate(res, "LogIn.html", nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) NotFound(res http.ResponseWriter, req *http.Request){
	if err := v.pagesTemplate.ExecuteTemplate(res, "Base.html", nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}