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

func (l *LocationServiceImpl) GetAllLocations() models.Result[[]string]{
	return l.Repository.GetAllLocations()
}

func (l *LocationServiceImpl) GetLocation(locationName string) models.Result[string]{
	return l.Repository.GetLocation(locationName)
}