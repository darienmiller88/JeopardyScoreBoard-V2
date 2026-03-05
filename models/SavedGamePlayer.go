package models

 type SavedGamePlayer struct {
	SavedGameID   int    `db:"saved_game_id"`
	PlayerID      int    `db:"player_id"`
	PlayerScore   int    `db:"player_score"`
	NameEncrypted []byte `db:"player_name_encrypted"`
	NameHash      []byte `db:"player_name_hash"`
}