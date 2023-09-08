// Package keycache contains an in-memory key-value database definition for storing user encryption keys
package keycache

import (
	"log"
	"sync"

	"github.com/grafviktor/keep-my-secret/internal/constant"
)

var (
	singleton *dataKeyCache
	once      sync.Once
)

type dataKeyCache struct {
	keymap map[string]string
}

// GetInstance - creates new cache storage for user data encryption keys.
// This is in-memory key-value storage
func GetInstance() *dataKeyCache {
	once.Do(
		func() {
			singleton = &dataKeyCache{
				keymap: make(map[string]string),
			}
		})

	return singleton
}

// Set - sets encryption key for a login name
func (u *dataKeyCache) Set(login, key string) {
	log.Printf("Set data key for user %s\n", login)

	u.keymap[login] = key
}

// Get - gets encryption key for a login name from the storage
func (u *dataKeyCache) Get(login string) (string, error) {
	key, ok := u.keymap[login]

	if !ok {
		log.Printf("Data key not found for user %s\n", login)

		return "", constant.ErrNotFound
	}

	return key, nil
}
