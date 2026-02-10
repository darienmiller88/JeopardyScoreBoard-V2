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
	UpdatePlayerName(oldPlayerName string, newPlayerName string, locationName string) models.Result[models.Player]
	AddPlayerToLocation(locationName string, player models.Player) models.Result[models.Player]
	GetPlayersFromLocation(locationName string) models.Result[[]models.Player]
	RemovePlayer(playerName string, locationName string) models.Result[models.Player]
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

	player.PlayerName = string(encryptedName)

	return utils.GetResult(nil, http.StatusCreated, player)
}

// Function to update a players name for a given location.
func (s *sqlPlayerRepository) UpdatePlayerName(oldPlayerName string, newPlayerName string, locationName string) models.Result[models.Player] {
	//Encrypt the new name
	encryptedName, err := s.encryptionService.Encrypt(newPlayerName)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	//Find the column of the old name by its hash, and create a new hash for the new name
	olderPlayerNameHash := utils.NameHash(oldPlayerName)
	newPlayerNameHash   := utils.NameHash(newPlayerName)

	//Use all of the following to updated the players name, encrypted name, and new hash.
	result, err := s.db.Exec(
		constants.UpdatePlayerName,
		newPlayerName, //$1 -> To be removed soon! 
		encryptedName, //$2
		newPlayerNameHash, //$3
		olderPlayerNameHash,//$3
		locationName, //$5
	)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	numRowsAffected, _ := result.RowsAffected()

	if numRowsAffected == 0 {
		return utils.GetResult(fmt.Errorf("could not find player %s", oldPlayerName), http.StatusNotFound, models.Player{})
	}

	updateResult :=  models.Player{
		PlayerName: newPlayerName,
		PlayerNameEncrypted: encryptedName,
		PlayerNameHash: newPlayerNameHash,
	}

	return utils.GetResult(nil, http.StatusOK, updateResult)
}

// Remove a single player from a given location.
func (s *sqlPlayerRepository) RemovePlayer(playerName string, locationName string) models.Result[models.Player] {
	result, err := s.db.Exec(constants.DeletePlayer, utils.NameHash(playerName), locationName)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Player{})
	}

	numRowsAffected, _ := result.RowsAffected()

	if numRowsAffected == 0 {
		return utils.GetResult(fmt.Errorf("could not find player %s", playerName), http.StatusNotFound, models.Player{})
	}

	return utils.GetResult(nil, http.StatusOK, models.Player{PlayerName: playerName})
}

// GetPlayersFromLocation fetches all players for a given location name,
// then decrypts their encrypted names before returning them to the caller.
//
// The database only stores encrypted names, so decryption must happen
// after retrieval and before returning API data.
func (s *sqlPlayerRepository) GetPlayersFromLocation(locationName string) models.Result[[]models.Player] {
	players := []models.Player{}

	// Query players by location (still contains encrypted names)
	if err := s.db.Select(&players, constants.GetAllPlayersFromLocation, locationName); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, players)
	}

	// Decrypt all player names before returning
	result := s.decryptPlayers(players)

	if result.Err != nil {
		return result
	}

	return utils.GetResult(nil, http.StatusOK, result.ResultData)
}


// GetAllPlayersFromAllLocations retrieves every player across all locations,
// then decrypts their names for safe API consumption.
func (s *sqlPlayerRepository) GetAllPlayersFromAllLocations() models.Result[[]models.Player] {
	players := []models.Player{}

	// Raw DB fetch (encrypted data)
	if err := s.db.Select(&players, constants.GetAllPlayers); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, players)
	}

	// Important: must return decrypted slice, not original
	result := s.decryptPlayers(players)

	if result.Err != nil {
		return result
	}

	return utils.GetResult(nil, http.StatusOK, result.ResultData)
}

// GetPlayersByNames retrieves players matching a list of plaintext names.
// The SQL uses hashes to locate the correct rows, then we decrypt
// the stored encrypted names before returning them.
func (s *sqlPlayerRepository) GetPlayersByNames(players []string) models.Result[[]models.Player] {
	validPlayers := []models.Player{}

	// pq.Array allows passing Go slice into SQL ANY() comparison
	if err := s.db.Select(&validPlayers, constants.GetPlayersByNames, pq.Array(players)); err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, validPlayers)
	}

	// Convert encrypted names into readable form
	result := s.decryptPlayers(validPlayers)

	if result.Err != nil {
		return result
	}

	return utils.GetResult(nil, http.StatusOK, result.ResultData)
}

// GetPlayerByName retrieves a single player by plaintext name.
// The lookup is performed via hash in SQL, then the encrypted
// name is decrypted before returning.
func (s *sqlPlayerRepository) GetPlayerByName(playerName string) models.Result[models.Player] {
	player := models.Player{}

	// Fetch the row (contains encrypted name)
	if err := s.db.Get(&player, constants.GetPlayerByName, playerName); err != nil {
		if err == sql.ErrNoRows {
			return utils.GetResult(err, http.StatusNotFound, player)
		}

		return utils.GetResult(err, http.StatusInternalServerError, player)
	}

	// Reuse bulk decrypt logic by wrapping in slice
	result := s.decryptPlayers([]models.Player{player})

	if result.Err != nil {
		return utils.GetResult(result.Err, result.StatusCode, result.ResultData[0])
	}

	// Return the decrypted version
	return utils.GetResult(nil, http.StatusOK, result.ResultData[0])
}

// decryptPlayers iterates over a slice of players and replaces the
// encrypted name with its decrypted equivalent in PlayerNameDecrypted.
//
// We pass the slice by value, but modify elements by index so the
// changes persist in the returned slice.
func (s *sqlPlayerRepository) decryptPlayers(players []models.Player) models.Result[[]models.Player] {
	for i := range players {

		// Decrypt the stored ciphertext from the DB
		decryptedName, err := s.encryptionService.Decrypt(players[i].PlayerNameEncrypted)
		if err != nil {
			return utils.GetResult(err, http.StatusInternalServerError, []models.Player{})
		}

		// Store decrypted value in the JSON-visible field
		players[i].PlayerNameDecrypted = decryptedName
	}

	return utils.GetResult(nil, http.StatusOK, players)
}
