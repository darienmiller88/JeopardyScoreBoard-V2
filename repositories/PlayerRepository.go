package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

const(
	push string = "$push"
	pull string = "$pull"
)

type PlayerRepository interface {
	UpdatePlayerName(oldPlayerName string, newPlayerName string)   models.Result[models.Player]
	AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player]
	GetPlayersFromLocation(locationName string)                    models.Result[[]models.Player]
	RemovePlayer(playerName string)                                models.Result[models.Player]
	GetAllPlayersFromAllLocations()                                models.Result[[]models.Player]
}

type sqlPlayerRepository struct{
	db *sqlx.DB
}

//Receive new Instance of MongoPlayerCardRepository.
func GetSqlPlayerCardRepository(newDB *sqlx.DB) *sqlPlayerRepository{
	return &sqlPlayerRepository{ db: newDB }
}

//Add a single player to a given location.
func (s *sqlPlayerRepository) AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player]{
	statement, err := s.db.PrepareNamed(constants.InsertNewPlayerWithoutTeam)

	if err != nil {
		return getResult(err, http.StatusInternalServerError, models.Player{})
	}

	params := map[string]interface{}{
		"player_name":  player.PlayerName,
		"location_name": locationName,
	}

	if err := statement.Get(&player.ID, params); err != nil{
		if err == sql.ErrNoRows {
			return getResult(fmt.Errorf("No location %s found", locationName), http.StatusNotFound, models.Player{})
		}
		
		return getResult(err, http.StatusInternalServerError, models.Player{})	
	}
	
	return getResult(nil, http.StatusOK, player)
}

//Function to update a players name for a given location.
func (s *sqlPlayerRepository) UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]{
	result, err := s.db.Exec(constants.UpdatePlayerName, newPlayerName, oldPlayerName)
	numRowsAffected, _ := result.RowsAffected()
	
	if numRowsAffected == 0 {
		return getResult(fmt.Errorf("could not find player %s", oldPlayerName), http.StatusNotFound, models.Player{})
	}
	
	if err != nil {	
		return getResult(err, http.StatusInternalServerError, models.Player{})
	}

	return getResult(nil, http.StatusOK, models.Player{ PlayerName: newPlayerName })
}

//Remove a single player from a given location.
func (s *sqlPlayerRepository) RemovePlayer(playerName string) models.Result[models.Player]{
	return getResult(nil, http.StatusOK, models.Player{})
}

func (s *sqlPlayerRepository) GetPlayersFromLocation(locationName string) models.Result[[]models.Player]{
	players := []models.Player{}

	return getResult(nil, http.StatusOK, players)
}

func (s *sqlPlayerRepository) GetAllPlayersFromAllLocations() models.Result[[]models.Player]{
	players := []models.Player{}

	if err := s.db.Select(&players, constants.GetAllPlayers); err != nil {
		return getResult(err, http.StatusInternalServerError, players)
	} 

	return getResult(nil, http.StatusOK, players)
}