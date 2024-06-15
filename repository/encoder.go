package repository

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"io"

	"golang.org/x/crypto/scrypt"
)

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

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	result := append(salt, ciphertext...)

	return result, nil
}

func decode(ciphertext []byte, password string) ([]Entry, error) {
	salt := ciphertext[:32]
	ciphertext = ciphertext[32:]

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
