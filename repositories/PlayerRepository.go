package repositories

import (
	"JeopardyScoreBoardV2/constants"
	"JeopardyScoreBoardV2/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
func GetSqlPlayerRepository(newDB *sqlx.DB) *sqlPlayerRepository{
	return &sqlPlayerRepository{ db: newDB }
}

//Add a single player to a given location.
func (s *sqlPlayerRepository) AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player]{
	if err := s.db.Get(&player.ID, constants.InsertNewPlayerWithoutTeam, player.PlayerName, locationName); err != nil{
		var pqErr *pq.Error
		
		if errors.As(err, &pqErr) {
            switch pqErr.Code {
            case "23502", "23503": // NOT NULL or FK violation
                return getResult(fmt.Errorf("no location '%s' found", locationName), http.StatusNotFound, models.Player{})
			case "23505"://unique key violation
				return getResult(fmt.Errorf("name %s is already taken", player.PlayerName), http.StatusConflict, models.Player{})
			}
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

	if err := s.db.Select(&players, constants.GetAllPlayersFromLocation, locationName); err != nil {
		return getResult(err, http.StatusInternalServerError, players)
	} 

	return getResult(nil, http.StatusOK, players)
}

func (s *sqlPlayerRepository) GetAllPlayersFromAllLocations() models.Result[[]models.Player]{
	players := []models.Player{}

	if err := s.db.Select(&players, constants.GetAllPlayers); err != nil {
		return getResult(err, http.StatusInternalServerError, players)
	} 

	return getResult(nil, http.StatusOK, players)
}