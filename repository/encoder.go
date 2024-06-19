package repository

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

var ErrInvalidPassword = errors.New("invalid password")

func encode(entries []Entry, password string) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(entries); err != nil {
		return nil, err
	}
	data := buf.Bytes()

	key, salt, err := deriveKey(password)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// generate password hash
	passwordHash, saltForPasswordHash, err := generatePasswordHash(password)

	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// File format:
	// 00:32 bytes salt for password hash
	// 32:64 bytes salt for passwordHash
	// 64:96 bytes for dataEncryptionSalt
	// 96:... encrypted data
	var result []byte

	result = append(result, saltForPasswordHash...)
	result = append(result, passwordHash...)
	result = append(result, salt...)
	result = append(result, ciphertext...)

	return result, nil
}

func decode(ciphertext []byte, password string) ([]Entry, error) {
	saltForPasswordHash := ciphertext[:32]
	passwordHash := ciphertext[32:64]

	if !verifyPassword(password, passwordHash, saltForPasswordHash) {
		return nil, ErrInvalidPassword
	}
	salt := ciphertext[64:96]
	ciphertext = ciphertext[96:]

	key, _, err := deriveKeyWithSalt(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	data, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	entries := []Entry{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func deriveKey(password string) ([]byte, []byte, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, err
	}
	key, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	return key, salt, err
}

func deriveKeyWithSalt(passphrase string, salt []byte) ([]byte, []byte, error) {
	key, err := scrypt.Key([]byte(passphrase), salt, 32768, 8, 1, 32)
	return key, salt, err
}

func generatePasswordHash(password string) ([]byte, []byte, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, err
	}
	passwordHash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	hash := sha256.Sum256(passwordHash)
	return hash[:], salt, err
}

func verifyPassword(password string, passwordHash, salt []byte) bool {
	key, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return false
	}
	newHash := sha256.Sum256(key)
	return hmac.Equal(newHash[:], passwordHash)
}
