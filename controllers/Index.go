package controllers

import (
	"github.com/go-chi/chi/v5"
)

type Index struct{
	Router *chi.Mux
}

func (i *Index) Init(){
	
}