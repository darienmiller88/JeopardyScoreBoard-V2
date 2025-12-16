package services

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
	"context"
) 

type LocationService struct{
	Repository repositories.LocationRepository
}

func (l *LocationService) GetAllLocations(ctx context.Context) models.Result[[]models.Location]{
	return l.Repository.GetAllLocations(ctx)
}

func (l *LocationService) GetLocation(ctx context.Context, locationName string) models.Result[models.Location]{
	return l.Repository.GetLocation(ctx, locationName)
}
//Add a new adapt location, which for now, I will not expose.
// func AddNewAdaptLocation(req *http.Request, location models.Location) models.Result[models.Location]{
// 	insertOneResult, err := database.GetLocationsCollection().InsertOne(req.Context(), location)
// 	locationResult := models.Result[models.Location]{}

// 	if err != nil{
// 		locationResult.Err = err
// 		locationResult.StatusCode = http.StatusInternalServerError

// 		return locationResult
// 	}
	
// 	location.ID = insertOneResult.InsertedID.(bson.ObjectID)
// 	locationResult.ResultData = location
// 	locationResult.StatusCode = http.StatusOK

// 	return locationResult
// }


//Retrieve one location from MongoDB.
