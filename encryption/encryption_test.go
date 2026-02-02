package encryption

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestService() *EncryptionService {
	// 32 byte key for AES-256
	key := []byte("12345678901234567890123456789012")
	return NewService(key)
}

func TestEncryptionWithName_Happy(t *testing.T) {
	svc := newTestService()
	nameToEncrypt := "Alice Smith"

	cipher, err := svc.Encrypt(nameToEncrypt)

	require.NoError(t, err)

	if cipher == "" {
		t.Fatal("expected ciphertext, got empty string")
	}

	//assert the name being encrypted and the encrypted name are not the same 
	require.NotEqual(t, nameToEncrypt, cipher)
}

func TestEncrypt_SamePlaintextDifferentCiphertext(t *testing.T) {
	svc := newTestService()

	c1, _ := svc.Encrypt("ALICE SMITH")
	c2, _ := svc.Encrypt("ALICE SMITH")

	if c1 == c2 {
		t.Fatal("expected different ciphertexts due to random nonce")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	svc := newTestService()

	original := "BOB JONES"

	cipher, err := svc.Encrypt(original)

	require.NoError(t, err)

	plain, err := svc.Decrypt(cipher)

	require.NoError(t, err)
	require.Equal(t, original, plain)
}

func TestEncrypt_BadKeySize(t *testing.T) {
	badSvc := NewService([]byte("short-key"))

	_, err := badSvc.Encrypt("test")

	require.Error(t, err)
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	svc := newTestService()

	cipher, _ := svc.Encrypt("CHARLIE")
	data, _ := base64.StdEncoding.DecodeString(cipher)

	// flip a byte
	data[len(data)-1] ^= 0xFF

	tampered := base64.StdEncoding.EncodeToString(data)

	_, err := svc.Decrypt(tampered)
	require.Error(t, err)
}

func TestDecrypt_WrongKey(t *testing.T) {
	svc1 := newTestService()
	cipher, _ := svc1.Encrypt("DAVID")

	// different key
	svc2 := NewService([]byte("aaaadaaabaaabaaaaaaaaaaakaaa/aaa"))

	_, err := svc2.Decrypt(cipher)
	require.Error(t, err)
}

func TestDecrypt_CiphertextTooShort(t *testing.T) {
	svc := newTestService()
	short := base64.StdEncoding.EncodeToString([]byte("tiny"))

	_, err := svc.Decrypt(short)
	if err == nil {
		t.Fatal("expected error due to short ciphertext")
	}
}