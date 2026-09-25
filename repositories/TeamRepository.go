package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/encryption"
	"JeopardyScoreBoardV2/models"
	"JeopardyScoreBoardV2/utils"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TeamRepository interface {

	//Get a team with all of the players on that team
	GetTeamWithAllPlayersDB(teamId int) models.Result[models.Team]

	// Get all team names (to be put on a select tag on the front end)
	GetAllTeamNamesDB()                 models.Result[[]string]

	// Get all teams in the database
	GetAllTeamsDB()                     models.Result[[]models.Team]

	// Get all teams by their IDs
	GetAllTeamsByIds(teamIds []int)     models.Result[[]models.Team]
}   

type teamRepository struct {
	db *sqlx.DB
	encryptionService *encryption.EncryptionService
}

// Receive new Instance of MongoPlayerCardRepository.
func NewTeamRepository(newDB *sqlx.DB, encryptionService *encryption.EncryptionService) TeamRepository {
	return &teamRepository{
		db: newDB, 
		encryptionService: encryptionService,
	}
}

//Get a team with all of the players on that team
func (s *teamRepository) GetTeamWithAllPlayersDB(teamId int) models.Result[models.Team]{
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
func (s *teamRepository) GetAllTeamNamesDB() models.Result[[]string]{
	teamNames := []string{}

	if err := s.db.Select(&teamNames, constants.GetAllTeamsByName); err != nil{
		return utils.GetResult(err, http.StatusInternalServerError, []string{})
	}	

	return utils.GetResult(nil, http.StatusOK, teamNames)
}

//Checks if a winning team exists
func (s *teamRepository) GetAllTeamsDB() models.Result[[]models.Team]{
	teams := []models.Team{}

	if err := s.db.Get(&teams, constants.GetAllTeams); err != nil{
		return utils.GetResult(err, http.StatusInternalServerError, []models.Team{})
	}

	return utils.GetResult(nil, http.StatusOK, teams)
}

//Checks if a winning team exists
func (s *teamRepository) GetAllTeamsByIds(teamIds []int) models.Result[[]models.Team]{
	teams := []models.Team{}

	if err := s.db.Get(&teams, constants.GetTeamsByIds, pq.Array(teamIds)); err != nil{
		return utils.GetResult(err, http.StatusInternalServerError, []models.Team{})
	}

	return utils.GetResult(nil, http.StatusOK, teams)
}