package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"JeopardyScoreBoardV2/services"

	"github.com/go-chi/chi/v5"
)

type ViewsController struct {
	templates        map[string]*template.Template
	Router           *chi.Mux
	SavedGameService services.SaveGameService
	LocationService  services.LocationService
	PlayerService    services.PlayerService
	TeamService      services.TeamService
}

func (v *ViewsController) Init(
	LocationService services.LocationService,
	PlayerService services.PlayerService,
	SavedGameService services.SaveGameService,
	TeamService services.TeamService,
) {
	v.Router = chi.NewRouter()
	v.templates = make(map[string]*template.Template)
	v.LocationService = LocationService
	v.PlayerService = PlayerService
	v.SavedGameService = SavedGameService
	v.TeamService = TeamService

	v.InitTemplateMap()

	v.Router.Get("/", v.CreateGame)
	v.Router.Get("/team-mode", v.TeamMode)
	v.Router.Get("/add-player", v.AddPlayerPage)
	v.Router.Get("/view-games", v.ViewGames)
	v.Router.Post("/add-saved-game", v.SaveGame)
	v.Router.NotFound(v.NotFound)
}

func (v *ViewsController) InitTemplateMap() {
	partialFiles, err := filepath.Glob("./templates/partials/*.html")
	if err != nil {
		panic(fmt.Sprintf("Error loading partials: %v", err))
	}

	entries, err := os.ReadDir("./templates/pages")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		name, _ := strings.CutSuffix(entry.Name(), ".html")

		files := []string{"templates/Base.html"}
		files = append(files, partialFiles...)
		files = append(files, fmt.Sprintf("templates/pages/%s.html", name))

		v.templates[name] = template.Must(template.ParseFiles(files...))
	}
}

func (v *ViewsController) SaveGame(res http.ResponseWriter, req *http.Request) {

}

func (v *ViewsController) CreateGame(res http.ResponseWriter, req *http.Request) {
	if err := v.templates["CreateGame"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) TeamMode(res http.ResponseWriter, req *http.Request) {
	teamNamesResult := v.TeamService.GetAllTeamNames()

	if teamNamesResult.Err != nil {
		http.Error(res, teamNamesResult.Err.Error(), teamNamesResult.StatusCode)
		return
	}

	//provide the team names to the template
	data := map[string]any{
		"TeamNames": teamNamesResult.ResultData,
	}

	if err := v.templates["TeamMode"].Execute(res, data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) AddPlayerPage(res http.ResponseWriter, req *http.Request) {
	locationsResult := v.LocationService.GetAllLocations()

	if locationsResult.Err != nil {
		http.Error(res, locationsResult.Err.Error(), locationsResult.StatusCode)
		return
	}

	selectedLocation := locationsResult.ResultData[0]
	playersResult := v.PlayerService.GetPlayersFromLocation(selectedLocation)

	if playersResult.Err != nil {
		http.Error(res, playersResult.Err.Error(), playersResult.StatusCode)
		return
	}

	data := map[string]any{
		"Players":          playersResult.ResultData,
		"Locations":        locationsResult.ResultData,
		"HasPlayers":       len(playersResult.ResultData) > 0,
		"SelectedLocation": selectedLocation,
	}

	if err := v.templates["AddPlayer"].Execute(res, data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) ViewGames(res http.ResponseWriter, req *http.Request) {
	savedGamesResult := v.SavedGameService.GetAllSavedGames()

	if savedGamesResult.Err != nil {
		http.Error(res, savedGamesResult.Err.Error(), http.StatusInternalServerError)
		return
	}

	data := template.FuncMap{
		"Games": savedGamesResult.ResultData,
	}

	if err := v.templates["ViewGames"].Execute(res, data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) NotFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	if err := v.templates["NotFound"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
