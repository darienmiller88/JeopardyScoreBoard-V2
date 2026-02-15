package encryption

import (
	"crypto/sha256"
	"os"
	"strings"
)

//Returns a 32 byte hash combining the player's name, and the pepper for the hashing
func NameHash(name string) []byte {
    pepper := os.Getenv("HASH_PEPPER") // different from encryption key
    h := sha256.Sum256([]byte(strings.TrimSpace(strings.ToLower(name)) + pepper))
    
	return h[:]
}