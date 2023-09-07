package middleware

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafviktor/keep-my-secret/internal/config"
)

func TestNewMiddleware(t *testing.T) {
	// Create a sample AppConfig for testing
	appConfig := config.AppConfig{}

	// Call the New function to create a middleware instance
	mw := New(appConfig)

	// Check if the config field of the middleware matches the expected AppConfig
	if mw.config != appConfig {
		t.Errorf("Expected middleware config to be %+v, but got %+v", appConfig, mw.config)
	}

	// Check if the config field of the middleware matches the expected AppConfig
	require.NotNil(t, mw.authVerifier)
}
