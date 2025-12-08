package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ViewsController struct{
	templates map[string]*template.Template
	Router    *chi.Mux
}

func (v *ViewsController) Init(){
	v.Router    = chi.NewRouter()
	v.templates = make(map[string]*template.Template)

	//Initialize template map
	v.InitTemplateMap()

	//Initialize view routes
	v.Router.Get("/", v.CreateGame)
	v.Router.Get("/team-mode", v.TeamMode)
	v.Router.Get("/add-player", v.AddPlayer)
	v.Router.Get("/view-games", v.ViewGames)
	v.Router.Get("/log-in", v.LogIn)
	v.Router.NotFound(v.NotFound)
}

func (v *ViewsController) InitTemplateMap(){
	entries, err := os.ReadDir("./templates/pages")

	if err != nil {
		panic(err)
	}else{
		for _, entry := range entries{
			name, _ := strings.CutSuffix(entry.Name(), ".html")
			v.templates[name] = template.Must(template.ParseFiles("templates/Base.html", fmt.Sprintf("templates/pages/%s.html", name)))
		}
	}
}

func (v *ViewsController) CreateGame(res http.ResponseWriter, req *http.Request){
	if err := v.templates["CreateGame"].Execute(res, nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) TeamMode(res http.ResponseWriter, req *http.Request){
	if err := v.templates["TeamMode"].Execute(res, nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) AddPlayer(res http.ResponseWriter, req *http.Request){
	if err := v.templates["AddPlayer"].Execute(res, nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) ViewGames(res http.ResponseWriter, req *http.Request){
	if err := v.templates["ViewGames"].Execute(res, nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) LogIn(res http.ResponseWriter, req *http.Request){
	if err := v.templates["LogIn"].Execute(res, nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) NotFound(res http.ResponseWriter, req *http.Request){
	if err := v.templates["NotFound"].Execute(res, nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}	
}