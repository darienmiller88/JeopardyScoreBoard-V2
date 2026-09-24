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

type teamService struct {
	TeamRepository repositories.TeamRepository
}

func NewTeamService(teamRepository repositories.TeamRepository) TeamService {
	return &teamService{
		TeamRepository: teamRepository,
	}
}

func (t *teamService) GetTeamWithAllPlayers(teamId int) models.Result[models.Team]{
	return t.TeamRepository.GetTeamWithAllPlayersDB(teamId)
}

func (t *teamService) GetAllTeamNames() models.Result[[]string]{
	return t.TeamRepository.GetAllTeamNamesDB()
}

func (t *teamService) GetAllTeams() models.Result[[]models.Team]{
	return t.TeamRepository.GetAllTeamsDB()
}