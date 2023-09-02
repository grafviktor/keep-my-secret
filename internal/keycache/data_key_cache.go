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

func GetInstance() *dataKeyCache {
	once.Do(
		func() {
			singleton = &dataKeyCache{
				keymap: make(map[string]string),
			}
		})

	return singleton
}

func (u *dataKeyCache) Set(login, key string) {
	log.Printf("Set data key for user %s\n", login)

	u.keymap[login] = key
}

func (u *dataKeyCache) Get(login string) (string, error) {
	key, ok := u.keymap[login]

	if !ok {
		log.Printf("Data key not found for user %s\n", login)

		return "", constant.ErrNotFound
	}

	return key, nil
}
