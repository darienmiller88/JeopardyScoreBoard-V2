package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"JeopardyScoreBoardV2/repositories"
	"JeopardyScoreBoardV2/services"
)

type Index struct{
	Router *chi.Mux
	viewsController ViewsController
	locationsController LocationsController
	playersController PlayersController
}

func (i *Index) InitControllers(db *sqlx.DB){
	i.Router = chi.NewRouter()

	//Initialize the views controller
	i.viewsController.Init()

	//Initialize the jeopardy controller, and chose the service and repo implementation
	i.locationsController.Init(&services.LocationServiceImpl{ Repository: repositories.GetSqlLocationRepository(db) })
	i.playersController.Init(&services.PlayerServiceImpl{ Repository: repositories.GetSqlPlayerCardRepository(db) })

	//Afterwards, mount the views router onto this router, which wiil be mounted onto the main chi router
	//in main.go
	i.Router.Mount("/", i.viewsController.Router)
	i.Router.Mount("/locations", i.locationsController.Router)
	i.Router.Mount("/players", i.playersController.Router)
}