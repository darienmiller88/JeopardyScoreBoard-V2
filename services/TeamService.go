package services

import (
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/repositories"
)

type TeamService interface {
	GetTeamWithAllPlayers(teamId int) models.Result[models.Team]
	GetAllTeamNames() models.Result[[]string]
	GetAllTeams() models.Result[[]models.Team]
}

type TeamServiceImpl struct {
	TeamRepository repositories.TeamRepository
}

func (t *TeamServiceImpl) GetTeamWithAllPlayers(teamId int) models.Result[models.Team]{
	return t.TeamRepository.GetTeamWithAllPlayersDB(teamId)
}

func (t *TeamServiceImpl) GetAllTeamNames() models.Result[[]string]{
	return t.TeamRepository.GetAllTeamNamesDB()
}

func (t *TeamServiceImpl) GetAllTeams() models.Result[[]models.Team]{
	return t.TeamRepository.GetAllTeamsDB()
}