package services

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
) 

type LocationService interface{
	
	// GetAllLocations retrieves all location names in a string slice
	GetAllLocations() models.Result[[]string]

	// GetLocation retrieves a single location name
	GetLocation(locationName string) models.Result[string]
}

type locationService struct{
	Repository repositories.LocationRepository
}

func NewLocationService(repo repositories.LocationRepository) LocationService{
	return &locationService{
		Repository: repo,
	}
}

//Retrieve all location name in a string slice
func (l *locationService) GetAllLocations() models.Result[[]string]{
	return l.Repository.GetAllLocations()
}

//Retrieve a single location name 
func (l *locationService) GetLocation(locationName string) models.Result[string]{
	return l.Repository.GetLocation(locationName)
}