package controllers

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ViewsController struct{
	Router *chi.Mux
	tmpl   *template.Template
}

func (v *ViewsController) Init(){
	v.Router = chi.NewRouter()
	v.tmpl   = template.Must(template.ParseFiles("templates/*.html", "templates/partials/*.html"))

	v.Router.Get("/", v.CreateGame)
}

func (v *ViewsController) CreateGame(res http.ResponseWriter, req *http.Request){
	if err := v.tmpl.ExecuteTemplate(res, "CreateGame", nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}