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

type PlayerRepository interface {
	UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player]
	AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player]
	GetPlayersFromLocation(locationName string) models.Result[[]models.Player]
	RemovePlayer(playerName string) models.Result[models.Player]
	GetAllPlayersFromAllLocations() models.Result[[]models.Player]
	GetPlayersByNames(players []string) models.Result[[]models.Player]
	GetPlayerByName(playerName string) models.Result[models.Player]
}

type sqlPlayerRepository struct {
	db                *sqlx.DB
	encryptionService *encryption.EncryptionService
}

// Receive new Instance of MongoPlayerCardRepository.
func GetSqlPlayerRepository(newDB *sqlx.DB, encryptionService *encryption.EncryptionService) *sqlPlayerRepository {
	return &sqlPlayerRepository{db: newDB, encryptionService: encryptionService}
}

// Add a single player to a given location.
func (s *sqlPlayerRepository) AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player] {	
	//First, encrypt the player name the client has sent, and receive the ciphertext.
	encryptedName, err := s.encryptionService.Encrypt(player.PlayerName)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	// //Calculate the hash from the player name, and some secret salt.
	hash := utils.NameHash(player.PlayerName)

	//Use the hash and encrypted names to create a new row, with the player name there temporarily.
	err = s.db.QueryRow(
		constants.InsertNewPlayerWithoutTeam,
		encryptedName,
		hash,
		player.PlayerName,
		locationName,
	).Scan(&player.ID, &player.CreatedAt, &player.UpdatedAt, &player.LocationID)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	player.PlayerName = encryptedName

	return utils.GetResult(nil, http.StatusCreated, player)
}

// Function to update a players name for a given location.
func (s *sqlPlayerRepository) UpdatePlayerName(oldPlayerName string, newPlayerName string) models.Result[models.Player] {
	//Encrypt the new name
	// encryptedName, err := s.encryptionService.Encrypt(newPlayerName)

	// if err != nil {
	// 	return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	// }

	result, err := s.db.Exec(constants.UpdatePlayerName, newPlayerName, oldPlayerName)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	numRowsAffected, _ := result.RowsAffected()

	if numRowsAffected == 0 {
		return utils.GetResult(fmt.Errorf("could not find player %s", oldPlayerName), http.StatusNotFound, models.Player{})
	}

	return utils.GetResult(nil, http.StatusOK, models.Player{PlayerName: newPlayerName})
}

// Remove a single player from a given location.
func (s *sqlPlayerRepository) RemovePlayer(playerName string) models.Result[models.Player] {
	result, err := s.db.Exec(constants.DeletePlayer, playerName)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	numRowsAffected, _ := result.RowsAffected()

	if numRowsAffected == 0 {
		return utils.GetResult(fmt.Errorf("could not find player %s", playerName), http.StatusNotFound, models.Player{})
	}

	return utils.GetResult(nil, http.StatusOK, models.Player{PlayerName: playerName})
}

func (s *sqlPlayerRepository) GetPlayersFromLocation(locationName string) models.Result[[]models.Player] {
	players := []models.Player{}

	if err := s.db.Select(&players, constants.GetAllPlayersFromLocation, locationName); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, players)
	}

	return utils.GetResult(nil, http.StatusOK, players)
}

func (s *sqlPlayerRepository) GetAllPlayersFromAllLocations() models.Result[[]models.Player] {
	players := []models.Player{}

	if err := s.db.Select(&players, constants.GetAllPlayers); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, players)
	}

	return utils.GetResult(nil, http.StatusOK, players)
}

func (s *sqlPlayerRepository) GetPlayersByNames(players []string) models.Result[[]models.Player] {
	validPlayers := []models.Player{}

	if err := s.db.Select(&validPlayers, constants.GetPlayersByNames, pq.Array(players)); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, validPlayers)
	}

	return utils.GetResult(nil, http.StatusOK, validPlayers)
}

func (s *sqlPlayerRepository) GetPlayerByName(playerName string) models.Result[models.Player] {
	player := models.Player{}

	if err := s.db.Get(&player, constants.GetPlayerByName, playerName); err != nil {
		if err == sql.ErrNoRows {
			return utils.GetResult(err, http.StatusNotFound, player)
		}

		return utils.GetResult(err, http.StatusInternalServerError, player)
	}

	return utils.GetResult(nil, http.StatusOK, player)
}
