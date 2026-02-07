package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

type EncryptionService struct {
	key []byte
}

// This key must be a 32 byte key
func NewService(key []byte) *EncryptionService {
	return &EncryptionService{key: key}
}

// Encrypt plaintext using AES-256-GCM
func (e *EncryptionService) Encrypt(plaintext string) ([]byte, error) {
	//Retrieve a new instance of a cipher algorithm wrapped in the GCM mode of transportation
	gcm, err := e.getNewGCM()

	if err != nil {
		return []byte{}, err
	}

	// Generate a random nonce (NEVER reuse nonces!). If the same nonce were used, the same word would
	//generate the same encrypted message. This is similar to hashing a password.
	nonce := make([]byte, gcm.NonceSize())

	//Read in a randomly generated number for the nonce (Number once)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt and authenticate. The nonce will be prepended to the encrypted message so it
	//can be decrypted later on
	// Format: nonce + ciphertext + auth tag
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	//Append the ciphertext to the nonce as one byte array
	return append(nonce, ciphertext...), nil
}

func (e *EncryptionService) Decrypt(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	gcm, err := e.getNewGCM()

	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (e *EncryptionService) getNewGCM() (cipher.AEAD, error) {
	//A block is the underlying algortihm responsible for scrambling plain text. By itself,
	//it just scrambles data without auth.
	block, err := aes.NewCipher(e.key)

	if err != nil {
		return nil, err
	}

	// GCM mode builds on top of the block cipher, and provides authenticated encryption
	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	return gcm, nil
}
