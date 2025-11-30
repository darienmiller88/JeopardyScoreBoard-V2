package controllers

import (
	"html/template"
	
	"github.com/go-chi/chi/v5"
)

type ViewsController struct{
	Router *chi.Mux
	tmpl *template.Template
}

func (v *ViewsController) Init(){
	v.tmpl = template.Must(template.ParseGlob("templates/*.html"))

}

func (v *ViewsController) CreateGame(){

}