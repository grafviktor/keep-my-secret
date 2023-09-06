package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	mathrand "math/rand"
	"time"
)

// IsUsernameConformsPolicy - checks if username conforms with security policy. Stub function.
// Normally we should have generic settings for username and password complexity. But I'm reluctant
// to complicate this logic for a pet project.
func IsUsernameConformsPolicy(username string) bool {
	return len(username) > 0
}

// IsPasswordConformsPolicy - checks if password conforms with security policy. Stub function.
func IsPasswordConformsPolicy(password string) bool {
	return len(password) > 0
}

var (
	aesKeyLength     = 24
	validRandomChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?~"
)

// GenerateRandomPassword - generates password which is used to encrypt user data
// Returns string of the same length as aesKeyLength
func GenerateRandomPassword() string {
	mathrand.Seed(time.Now().UnixNano())

	key := make([]byte, aesKeyLength)
	for i := range key {
		key[i] = validRandomChars[mathrand.Intn(len(validRandomChars))]
	}

	return string(key)
}

func normalizeAESKey(key string) string {
	length := len(key)

	if length < aesKeyLength {
		for i := 0; i < aesKeyLength-length; i++ {
			key += "0"
		}
	} else if length > aesKeyLength {
		key = key[:aesKeyLength]
	}

	return key
}

// Encrypt - encrypts data with AES using key. If key is less that supported AES key length,
// then it is padded with zeros.
// plaindata - data to be encrypted
// key - key to be used for encryption
// Returns encrypted data or error
func Encrypt(plaindata []byte, key string) ([]byte, error) {
	key = normalizeAESKey(key)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaindata))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaindata)

	return ciphertext, nil
}

// Decrypt - decrypts data with AES using key. If key is less that supported AES key length,
// then it is padded with zeros.
// cipherdata - data to be decrypted
// key - key to be used for decryption
// Returns decrypted data or error
func Decrypt(cipherdata []byte, key string) ([]byte, error) {
	key = normalizeAESKey(key)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	iv := cipherdata[:aes.BlockSize]
	cipherdata = cipherdata[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(cipherdata))
	cfb.XORKeyStream(plaintext, cipherdata)

	return plaintext, nil
}
