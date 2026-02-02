package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type EncryptionService struct {
	key []byte
}

//This key must be a 32 byte key
func NewService(key []byte) *EncryptionService {
	return &EncryptionService{key: key}
}

// Encrypt plaintext using AES-256-GCM
func (e *EncryptionService) Encrypt(plaintext string) (string, error) {
	//Retrieve a new instance of a cipher algorithm wrapped in the GCM mode of transportation
	gcm, err := e.getNewGCM()

	if err != nil {
		return "", err
	}

	// Generate a random nonce (NEVER reuse nonces!). If the same nonce were used, the same word would 
	//generate the same encrypted message. This is similar to hashing a password.
	nonce := make([]byte, gcm.NonceSize())
	
	//Read in a randomly generated number for the nonce (Number once)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and authenticate. The nonce will be prepended to the encrypted message so it 
	//can be decrypted later on
	// Format: nonce + ciphertext + auth tag
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64 for text storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt ciphertext
func (e *EncryptionService) Decrypt(ciphertext string) (string, error) {
    if ciphertext == "" {
        return "", nil
    }

    // Decode from string into a byte array 
    data, err := base64.StdEncoding.DecodeString(ciphertext)

	if err != nil {
		return "", err
	}

	//Get a new instance of a GCM (Galois/Counter mode) cipher algorithm
    gcm, err := e.getNewGCM()

	if err != nil {
		return "", err
	}

    nonceSize := gcm.NonceSize()
	
	//Check to see if the length of the byte array converted from the ciphertext is less than the expected
	//nonce size. If so, return early to not unnecessarily decrypt.
    if len(data) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    // Split nonce and ciphertext
    nonce, encryptedName := data[:nonceSize], data[nonceSize:]

    // Decrypt and verify
    plaintext, err := gcm.Open(nil, nonce, encryptedName, nil)
   
	if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

func (e *EncryptionService) getNewGCM() (cipher.AEAD, error){
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