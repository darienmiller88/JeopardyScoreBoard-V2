package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type TeamRepository interface {
	GetTeamWithAllPlayersDB(teamId int) models.Result[models.Team]
	GetAllTeamsWithAllPlayersDB()       models.Result[[]models.Team]
	GetAllTeamNames()                   models.Result[string]
}   

type sqlTeamRepository struct {
	db *sqlx.DB
}

// Receive new Instance of MongoPlayerCardRepository.
func GetSqlTeamRepository(newDB *sqlx.DB) *sqlTeamRepository {
	return &sqlTeamRepository{db: newDB}
}

func (s *sqlTeamRepository) GetTeamWithAllPlayersDB(teamId int) models.Result[models.Team]{
	team := models.Team{}

	if err := s.db.Get(&team, constants.GetTeamById, teamId); err != nil{
		if err == sql.ErrNoRows {
			return getResult(fmt.Errorf("No team found with id %d", teamId), http.StatusNotFound, team)	
		} 

		return getResult(err, http.StatusInternalServerError, team)
	}

	players := []models.Player{}

	if err := s.db.Select(&players, constants.GetAllPlayersOnTeam, teamId); err != nil{
		return getResult(err, http.StatusInternalServerError, team)
	}
	
	team.Players = players

	return getResult(nil, http.StatusOK, team)
}	

func (s *sqlTeamRepository) GetAllTeamsWithAllPlayersDB() models.Result[[]models.Team]{
	return getResult(nil, http.StatusOK, []models.Team{})
}

func (s *sqlTeamRepository) GetAllTeamNames() models.Result[string]{
	return getResult(nil, http.StatusOK, "")
}
