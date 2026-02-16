package services

import (
	"JeopardyScoreBoardV2/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockTeamRepository struct {
	GetTeamWithAllPlayersDBFunc func(teamId int) models.Result[models.Team]
	GetAllTeamNamesDBFunc       func() models.Result[[]string]
	GetAllTeamsDBFunc           func() models.Result[[]models.Team]
	GetAllTeamsByIdsDBFunc      func(teamIds []int) models.Result[[]models.Team]
}

func (m *MockTeamRepository) GetTeamWithAllPlayersDB(teamId int) models.Result[models.Team] {
	return m.GetTeamWithAllPlayersDBFunc(teamId)
}

func (m *MockTeamRepository) GetAllTeamNamesDB() models.Result[[]string] {
	return m.GetAllTeamNamesDBFunc()
}

func (m *MockTeamRepository) GetAllTeamsDB() models.Result[[]models.Team] {
	return m.GetAllTeamsDBFunc()
}

func (m *MockTeamRepository) GetAllTeamsByIds(teamIds []int) models.Result[[]models.Team] {
	return m.GetAllTeamsByIdsDBFunc(teamIds)
}


func TestGetTeamWithAllPlayers_Happy(t *testing.T) {
	id := 1
	mockRepo := &MockTeamRepository{
		GetTeamWithAllPlayersDBFunc: func(teamId int) models.Result[models.Team] {
			return models.Result[models.Team]{
				ResultData: models.Team{ID: id},
				StatusCode: http.StatusOK,
			}
		},
	}
	service := &TeamServiceImpl{TeamRepository: mockRepo}
	result := service.GetTeamWithAllPlayers(id)

	assert.Equal(t, id, result.ResultData.ID)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestGetTeamWithAllPlayers_TeamNotFound_Unhappy(t *testing.T) {
	id := 765//id does not exist in DB
	mockRepo := &MockTeamRepository{
		GetTeamWithAllPlayersDBFunc: func(teamId int) models.Result[models.Team] {
			return models.Result[models.Team]{
				ResultData: models.Team{ID: id},
				StatusCode: http.StatusNotFound,
			}
		},
	}
	service := &TeamServiceImpl{TeamRepository: mockRepo}
	result := service.GetTeamWithAllPlayers(id)

	assert.Equal(t, id, result.ResultData.ID)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}

func TestGetAllTeamNames_Happy(t *testing.T) {
	teams := []string{"Team A", "Team B"}
	mockRepo := &MockTeamRepository{
		GetAllTeamNamesDBFunc: func() models.Result[[]string] {
			return models.Result[[]string]{
				ResultData: teams,
				StatusCode: http.StatusOK,
			}
		},
	}
	service := &TeamServiceImpl{TeamRepository: mockRepo}
	result := service.GetAllTeamNames()

	assert.Equal(t, len(teams), len(result.ResultData))
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestGetAllTeams_Happy(t *testing.T) {
	mockRepo := &MockTeamRepository{
		GetAllTeamsDBFunc: func() models.Result[[]models.Team] {
			return models.Result[[]models.Team]{
				ResultData: []models.Team{{ID: 1}, {ID: 2}},
				StatusCode: http.StatusOK,
			}
		},
	}
	service := &TeamServiceImpl{TeamRepository: mockRepo}
	result := service.GetAllTeams()

	assert.Equal(t, 2, len(result.ResultData))
	assert.Equal(t, http.StatusOK, result.StatusCode)
}
