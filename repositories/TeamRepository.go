package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/utils"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type TeamRepository interface {
	GetTeamWithAllPlayersDB(teamId int) models.Result[models.Team]
	GetAllTeamNames()                   models.Result[[]string]
}   

type sqlTeamRepository struct {
	db *sqlx.DB
}

// Receive new Instance of MongoPlayerCardRepository.
func GetSqlTeamRepository(newDB *sqlx.DB) *sqlTeamRepository {
	return &sqlTeamRepository{db: newDB}
}

//Get a team with all of the players on that team
func (s *sqlTeamRepository) GetTeamWithAllPlayersDB(teamId int) models.Result[models.Team]{
	team := models.Team{}

	if err := s.db.Get(&team, constants.GetTeamById, teamId); err != nil{
		if err == sql.ErrNoRows {
			return utils.GetResult(fmt.Errorf("No team found with id %d", teamId), http.StatusNotFound, team)	
		} 

		return utils.GetResult(err, http.StatusInternalServerError, team)
	}

	players := []models.Player{}

	if err := s.db.Select(&players, constants.GetAllPlayersOnTeam, teamId); err != nil{
		return utils.GetResult(err, http.StatusInternalServerError, team)
	}
	
	team.Players = players

	return utils.GetResult(nil, http.StatusOK, team)
}	

//Get all team names (to be put on a select tag on the front end)
func (s *sqlTeamRepository) GetAllTeamNames() models.Result[[]string]{
	teamNames := []string{}

	if err := s.db.Select(&teamNames, constants.GetAllTeamsByName); err != nil{
		return utils.GetResult(err, http.StatusInternalServerError, []string{})
	}	

	return utils.GetResult(nil, http.StatusOK, teamNames)
}
