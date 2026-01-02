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
	jeopardyController JeopardyController
}

func (i *Index) InitControllers(db *sqlx.DB){
	i.Router = chi.NewRouter()

	//Initialize the views controller
	i.viewsController.Init()

	//Initialize the jeopardy controller, and chose the service and repo implementation
	i.jeopardyController.Init(&services.LocationServiceImpl{ Repository: repositories.GetSqlLocationRepository(db) })

	//Afterwards, mount the views router onto this router, which wiil be mounted onto the main chi router
	//in main.go
	i.Router.Mount("/", i.viewsController.Router)
	i.Router.Mount("/locations", i.jeopardyController.Router)
}