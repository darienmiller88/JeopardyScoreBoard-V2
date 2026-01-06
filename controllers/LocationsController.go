package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"JeopardyScoreBoardV2/services"
)

type LocationsController struct{
	Router *chi.Mux
	locationService services.LocationService
}

func (l *LocationsController) Init(service services.LocationService){
	l.Router = chi.NewRouter()
	l.locationService = service

	l.Router.Get("/", l.getLocations)
	l.Router.Get("/{location_name}", l.getLocationByName)
}

func (l *LocationsController) getLocations(res http.ResponseWriter, req *http.Request){
	result := l.locationService.GetAllLocations()

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	data, err := json.Marshal(&result)

	if err != nil{
		http.Error(res, err.Error(), 500)
		return
	}
	
	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	res.Write(data)
}

func (l *LocationsController) getLocationByName(res http.ResponseWriter, req *http.Request){
	locationName := chi.URLParam(req, "location_name")
	result := l.locationService.GetLocation(locationName)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	data, err := json.Marshal(&result)

	if err != nil{
		http.Error(res, err.Error(), 500)
		return
	}
	
	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(200)
	res.Write(data)
}