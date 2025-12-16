package services

import (
	"context"

	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PlayerCardService struct{
	Repository repositories.PlayerCardRepository
}

func (p *PlayerCardService) UpdatePlayerName(ctx context.Context, locationName string, oldPlayerName string, newPlayerName string)models.Result[*mongo.UpdateResult]{
	return p.Repository.UpdatePlayerName(ctx, locationName, oldPlayerName, newPlayerName)
}

func (p *PlayerCardService) AddPlayerToLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult] {
	return p.Repository.AddPlayerToLocation(ctx, locationName, playerName)
}

func (p *PlayerCardService) RemovePlayerFromLocation(ctx context.Context, locationName string, playerName string) models.Result[*mongo.UpdateResult]{
	return  p.Repository.RemovePlayerFromLocation(ctx, locationName, playerName)
}