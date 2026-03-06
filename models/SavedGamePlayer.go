package models

import "time"

type SavedGamePlayer struct {
	ID            int       `db:"id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	PlayerID      int       `db:"player_id"`
	SavedGameID   int       `db:"saved_game_id"`
	PlayerScore   int       `db:"player_score"`
	NameEncrypted []byte    `db:"player_name_encrypted"`
	NameHash      []byte    `db:"player_name_hash"`
}