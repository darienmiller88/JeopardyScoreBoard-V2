package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"JeopardyScoreBoardV2/services"
)

type ViewsController struct{
	templates        map[string]*template.Template
	Router          *chi.Mux
	SavedGameService services.SaveGameService
	LocationService  services.LocationService
	PlayerService    services.PlayerService
}

func (v *ViewsController) Init(LocationService services.LocationService, PlayerService services.PlayerService){
	v.Router    = chi.NewRouter()
	v.templates = make(map[string]*template.Template)
	v.LocationService = LocationService
	v.PlayerService = PlayerService

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
	// Get all partial files
	partialFiles, err := filepath.Glob("./templates/partials/*.html")

	if err != nil {
		panic(fmt.Sprintf("Error loading partials: %v", err))
	}

	//Get all pages
	entries, err := os.ReadDir("./templates/pages")

	if err != nil {
		panic(err)
	}

	for _, entry := range entries{
		name, _ := strings.CutSuffix(entry.Name(), ".html")
		
		//For each page, build the following file stack: Base html, all partials, Page.html
		files := []string{"templates/Base.html"}
		files = append(files, partialFiles...)
		files = append(files, fmt.Sprintf("templates/pages/%s.html", name))
		
		v.templates[name] = template.Must(template.ParseFiles(files...))
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
	locationsResult := v.LocationService.GetAllLocations()

	if locationsResult.Err != nil {
		http.Error(res, locationsResult.Err.Error(), locationsResult.StatusCode)
		return
	}

	playersResult := v.PlayerService.GetPlayersFromLocation(locationsResult.ResultData[3])

	if playersResult.Err != nil {
		http.Error(res, playersResult.Err.Error(), playersResult.StatusCode)
		return
	}

	data := map[string]any{
		"Locations": locationsResult.ResultData,
		"Players": playersResult.ResultData,
	}

	if err := v.templates["AddPlayer"].Execute(res, data); err != nil{
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
	res.WriteHeader(http.StatusNotFound)
	if err := v.templates["NotFound"].Execute(res, nil); err != nil{
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}	
}