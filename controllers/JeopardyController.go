package controllers

import (
	"encoding/json"
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
	result := j.locationService.GetAllLocations()

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

func (j *JeopardyController) getLocationByName(res http.ResponseWriter, req *http.Request){
	locationName := chi.URLParam(req, "location_name")
	result := j.locationService.GetLocation(locationName)

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