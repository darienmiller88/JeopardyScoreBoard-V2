package utils

import (
	"JeopardyScoreBoardV2/models"
	"crypto/sha256"
	"os"
	"strings"
)

//Helper function to allow repos to send result payloads with less text.
func GetResult[T any](err error, statusCode int, payload T) models.Result[T] {
	return models.Result[T]{
		StatusCode: statusCode,
		Err: err,
		ResultData: payload,
	}
}

//Returns a 32 byte hash combining the player's name, and the salt for the hashing
func NameHash(name string) []byte {
    pepper := os.Getenv("HASH_SALT") // different from encryption key
    h := sha256.Sum256([]byte(strings.ToLower(name) + pepper))
    
	return h[:]
}
