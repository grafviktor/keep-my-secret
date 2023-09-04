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

type Encryptor interface {
	Encrypt(plaindata []byte, key string) ([]byte, error)
	Decrypt(cipherdata []byte, key string) ([]byte, error)
}

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

func (s *Secret) SetEncryptor(encryptor Encryptor) {
	s.Encryptor = encryptor
}

const (
	typeString = "string"
	typeBinary = "[]uint8"
)

func (s *Secret) Encrypt(key, salt string) error {
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

		var encrypted []byte
		var err error

		if s.Encryptor == nil {
			encrypted, err = utils.Encrypt(toEncrypt, key)
		} else {
			encrypted, err = s.Encryptor.Encrypt(toEncrypt, key)
		}

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

func (s *Secret) Decrypt(key, salt string) error {
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

		var decrypted []byte
		var err error

		if s.Encryptor == nil {
			decrypted, err = utils.Decrypt(toDecrypt, key)
		} else {
			decrypted, err = s.Encryptor.Decrypt(toDecrypt, key)
		}

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
