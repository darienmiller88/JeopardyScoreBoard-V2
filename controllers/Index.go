package controllers

import (
	"github.com/go-chi/chi/v5"
)

type Index struct{
	Router *chi.Mux
	viewsController ViewsController
}

func (i *Index) Init(){
	i.Router = chi.NewRouter()

	//Initialize the views controller
	i.viewsController.Init()

	//Afterwards, mount the views router onto this router, which wiil be mounted onto the main chi router
	//in main.go
	i.Router.Mount("/", i.viewsController.Router)
}