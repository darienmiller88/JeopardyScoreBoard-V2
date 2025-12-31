package controllers

import(
	"net/http"
	"github.com/go-chi/chi/v5"

	"JeopardyScoreBoardV2/services"
)

type JeopardyController struct{
	Router *chi.Mux
	locationService services.LocationService
}

func (j *JeopardyController) Init(service services.LocationService){
	j.Router = chi.NewRouter()
	j.locationService = service

	j.Router.Get("/", j.getLocations)
	j.Router.Get("/{location_name}", j.getLocationByName)
}

func (j *JeopardyController) getLocations(res http.ResponseWriter, req *http.Request){
	

}

func (j *JeopardyController) getLocationByName(res http.ResponseWriter, req *http.Request){
	
	
}