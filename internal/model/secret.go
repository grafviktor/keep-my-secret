package model

import (
	"fmt"
	"log"
	"reflect"

	"github.com/samber/lo"

	"github.com/grafviktor/keep-my-secret/internal/api/utils"
)

// var shouldNotEncrypt = []string{"ID", "Type", "Title"}
var shouldNotEncrypt = []string{"ID", "Encryptor"}

// Encryptor is used for setting encrypting method for Secret model. This interface is used mainly for mocking
type Encryptor interface {
	Encrypt(secret *Secret, key, salt string) error
	Decrypt(secret *Secret, key, salt string) error
}

// Secret is a model of secret object which the application receives from the client
type Secret struct {
	ID             int64     `json:"id"`
	Type           string    `json:"type"`
	Title          string    `json:"title"`
	Login          string    `json:"login"`
	Password       string    `json:"password"`
	Note           string    `json:"note"`
	File           []byte    `json:"-"`
	FileName       string    `json:"file_name"`
	CardholderName string    `json:"cardholder_name"`
	CardNumber     string    `json:"card_number"`
	Expiration     string    `json:"expiration"`
	SecurityCode   string    `json:"security_code"`
	Encryptor      Encryptor `json:"-"`
}

// SetEncryptor should be used for setting concrete encryptor implementation. Currently used in unit tests
func (s *Secret) SetEncryptor(encryptor Encryptor) {
	s.Encryptor = encryptor
}

const (
	typeString = "string"
	typeBinary = "[]uint8"
)

// Encrypt - encrypts object using key and salt
func (s *Secret) Encrypt(key, salt string) error {
	if s.Encryptor != nil {
		return s.Encryptor.Encrypt(s, key, salt)
	}

	v := reflect.Indirect(reflect.ValueOf(s))
	typeOfP := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := typeOfP.Field(i).Name

		// skip fields that should not be decrypted
		if _, ok := lo.Find(shouldNotEncrypt, func(i string) bool {
			return i == fieldName
		}); ok {
			continue
		}

		fieldType := field.Type().String()
		var toEncrypt []byte

		switch {
		case fieldType == typeString:
			fieldValue, _ := field.Interface().(string)
			toEncrypt = []byte(salt + fieldValue)
		case fieldType == typeBinary:
			fieldValue, _ := field.Interface().([]byte)
			toEncrypt = append([]byte(salt), fieldValue...)
		default:
			log.Printf("secret decrypt: field %s is not a string", fieldName)
			continue
		}

		encrypted, err := utils.Encrypt(toEncrypt, key)
		if err != nil {
			return fmt.Errorf("secret.Encrypt: %s", err.Error())
		}

		if fieldType == typeString {
			v.Field(i).SetString(string(encrypted))
		} else {
			v.Field(i).SetBytes(encrypted)
		}
	}

	return nil
}

// Decrypt - decrypts object using key and salt
func (s *Secret) Decrypt(key, salt string) error {
	if s.Encryptor != nil {
		return s.Encryptor.Decrypt(s, key, salt)
	}

	v := reflect.Indirect(reflect.ValueOf(s))
	typeOfP := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := typeOfP.Field(i).Name

		// skip fields that should not be decrypted
		if _, ok := lo.Find(shouldNotEncrypt, func(i string) bool {
			return i == fieldName
		}); ok {
			continue
		}

		fieldType := field.Type().String()
		var toDecrypt []byte

		switch {
		case fieldType == typeString:
			fieldValue, _ := field.Interface().(string)
			toDecrypt = []byte(fieldValue)
		case fieldType == typeBinary:
			fieldValue, _ := field.Interface().([]byte)
			toDecrypt = fieldValue
		default:
			log.Printf("secret.Decrypt: field %s is not a string", fieldName)
			continue
		}

		decrypted, err := utils.Decrypt(toDecrypt, key)
		if err != nil {
			return fmt.Errorf("secret encrypt: %s", err.Error())
		}

		if fieldType == typeString {
			decryptedStr := ""
			if len(decrypted) > len(salt) {
				decryptedStr = string(decrypted[len(salt):])
			}

			v.Field(i).SetString(decryptedStr)
		} else {
			decryptedBytes := []byte{}
			if len(decrypted) > len(salt) {
				decryptedBytes = decrypted[len(salt):]
			}

			v.Field(i).SetBytes(decryptedBytes)
		}
	}

	return nil
}
