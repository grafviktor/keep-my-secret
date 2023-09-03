package keycache

import (
	"errors"
	"testing"

	"github.com/grafviktor/keep-my-secret/internal/constant"
)

func TestSingletonDataKeyCache(t *testing.T) {
	// Create a new singleton instance
	cache := GetInstance()

	// Test setting and getting data key
	login := "testuser"
	key := "testkey"

	cache.Set(login, key)

	result, err := cache.Get(login)
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}

	if result != key {
		t.Errorf("Expected data key %s, but got %s", key, result)
	}

	// Test getting a non-existent data key
	nonExistentLogin := "nonexistentuser"
	_, err = cache.Get(nonExistentLogin)

	if !errors.Is(err, constant.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, but got %v", err)
	}
}
