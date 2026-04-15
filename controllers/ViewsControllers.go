package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"JeopardyScoreBoardV2/middlewares"
	"JeopardyScoreBoardV2/services"

	"github.com/go-chi/chi/v5"
)

// Hardcoded users
var allowedUsers = map[string]bool{
	"jennamandelricci": true,
	"asisatmuldoon":    true,
	"midrenelamy":      true,
	"adonisbrown":      true,
	"lindalaul":        true,
}

var passwordHash = []byte("$2a$12$9zqH0OZV6FZCqzQ6lK3Q3u6F9mQ9pYJx0kZcZ5Y8QeV1wF6ZrY9eK")

type ViewsController struct {
	templates        map[string]*template.Template
	logInTemplate    *template.Template
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
	v.logInTemplate = template.Must(template.ParseFiles("templates/LogIn.html"))

	//Initialize template map
	v.InitTemplateMap()

	v.Router.Use(middlewares.AuthMiddleware)

	//Initialize view routes
	v.Router.Get("/", v.CreateGame)
	v.Router.Get("/team-mode", v.TeamMode)
	v.Router.Get("/add-player", v.AddPlayerPage)
	v.Router.Get("/view-games", v.ViewGames)
	v.Router.Get("/talent-show", v.TalentShow)
	v.Router.Get("/log-in", v.LogIn)
	v.Router.Get("/sign-out", v.SignOut)
	v.Router.Post("/log-in", v.HandleLogIn)
	v.Router.NotFound(v.NotFound)
}

func (v *ViewsController) InitTemplateMap() {
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

	for _, entry := range entries {
		name, _ := strings.CutSuffix(entry.Name(), ".html")

		//For each page, build the following file stack: Base html, all partials, Page.html
		files := []string{"templates/Base.html"}
		files = append(files, partialFiles...)
		files = append(files, fmt.Sprintf("templates/pages/%s.html", name))

		v.templates[name] = template.Must(template.ParseFiles(files...))
	}
}

func (v *ViewsController) TalentShow(res http.ResponseWriter, req *http.Request) {
	data := map[string]any{
		"TalentShowSlots": []string{
			"Christopher Taylor",
			"Kiefer Inson",
			"Tony Switzer",
			"CAYENNE NO_LUCK aka Justin Jacob",
			"Jadel Nunez",
			"Carla O’Brien & Josh Wilson",
			"Chloe Crisano",
			"Tony B. Rivers",
			"Money",
			"Kenny Shiver & Angie Eason",
			"Rachel Fonseca & Sophie Thurschwell",
			"Carlos Mendoza",
			"Woody Tanor",
			"Denise Farmer",
		},
	}

	if err := v.templates["TalentShow"].Execute(res, data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
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

func (v *ViewsController) LogIn(res http.ResponseWriter, req *http.Request) {
	if err := v.logInTemplate.Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewsController) HandleLogIn(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	if !allowedUsers[username] {
		fmt.Println("Invalid username")
		http.Redirect(res, req, "/log-in", http.StatusSeeOther)
		return
	}

	fmt.Println("username:", username, "password:", password)
	correctPassword := os.Getenv("PASSWORD")

	fmt.Println("value:", correctPassword == password)
	if password != correctPassword {
		fmt.Println("Invalid password")
		http.Redirect(res, req, "/log-in", http.StatusSeeOther)
		return
	}

	// ✅ Set session cookie
	http.SetCookie(res, &http.Cookie{
		Name:     "auth",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(res, req, "/talent-show", http.StatusSeeOther)
}

func (v *ViewsController) SignOut(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // 🔥 deletes cookie
		HttpOnly: true,
	})

	http.Redirect(res, req, "/log-in", http.StatusSeeOther)
}

func (v *ViewsController) NotFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	if err := v.templates["NotFound"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
