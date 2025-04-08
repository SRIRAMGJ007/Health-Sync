package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// Generate a fixed 32-byte AES-256 key from a passphrase
func GenerateKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

// EncryptData encrypts the given data using AES-256-GCM
func EncryptData(plainData, key []byte) ([]byte, error) {
	// Ensure key is 32 bytes (AES-256 requires 256-bit key)
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes long")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create AES cipher: " + err.Error())
	}

	// Use GCM (Galois/Counter Mode) for authentication
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM: " + err.Error())
	}

	// Generate a random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("failed to generate nonce: " + err.Error())
	}

	// Encrypt and prepend nonce
	encryptedData := aesGCM.Seal(nonce, nonce, plainData, nil)
	return encryptedData, nil
}

// DecryptData decrypts AES-256-GCM encrypted data
func DecryptData(encryptedData, key []byte) ([]byte, error) {
	// Ensure key is 32 bytes
	if len(key) != 32 {
		return nil, errors.New("decryption key must be 32 bytes long")
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create AES cipher: " + err.Error())
	}

	// Use GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM: " + err.Error())
	}

	// Extract nonce from the beginning of the encrypted data
	nonceSize := aesGCM.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("invalid encrypted data: missing nonce")
	}

	nonce, cipherText := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt data
	plainData, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt data: " + err.Error())
	}

	return plainData, nil
}
