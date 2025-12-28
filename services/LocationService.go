package services

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
) 

type LocationService struct{
	Repository repositories.LocationRepository
}

func (l *LocationService) GetAllLocations() models.Result[[]string]{
	return l.Repository.GetAllLocations()
}

func (l *LocationService) GetLocation(locationName string) models.Result[string]{
	return l.Repository.GetLocation(locationName)
}