package services

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
) 

type LocationService interface{
	GetAllLocations() models.Result[[]string]
	GetLocation(locationName string) models.Result[string]
}

type LocationServiceImpl struct{
	Repository repositories.LocationRepository
}

//Retrieve all location name in a string slice
func (l *LocationServiceImpl) GetAllLocations() models.Result[[]string]{
	return l.Repository.GetAllLocations()
}

//Retrieve a single location name 
func (l *LocationServiceImpl) GetLocation(locationName string) models.Result[string]{
	return l.Repository.GetLocation(locationName)
}