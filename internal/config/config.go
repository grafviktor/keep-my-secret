// Package config - contains application configuration structures
package config

import "github.com/grafviktor/keep-my-secret/internal/storage"

// EnvConfig is reqyured
type EnvConfig struct {
	// Secret which is used for signing cookies
	Secret string `env:"APP_SECRET"       envDefault:"romeo romeo whiskey"`
	// ServerAddr defines host and port where server is running
	ServerAddr string `env:"SERVER_ADDRESS"   envDefault:"localhost:8080"`
	// PSQL connection string (DSN)
	DSN string `env:"DSN"                     envDefault:"./kms.db"`
	// HTTPS TLS certificate path
	HTTPSCertPath string `env:"TLS_CERT_PATH" envDefault:"./tls/cert.pem"`
	// HTTPS TLS key path
	HTTPSKeyPath string `env:"TLS_KEY_PATH"   envDefault:"./tls/key.pem"`
	// Self-explanatory
	Domain string `env:"DOMAIN"               envDefault:"localhost"`
	// ClientURL is used to define browser path to the client application
	ClientURL string `env:"CLIENT_URL"        envDefault:"/"`
	// DevMode enables CORS
	DevMode bool `env:"DEV"                   envDefault:"false"`
}

type AppConfig struct {
	ServerAddr string
	// PSQL connection string (DSN)
	DSN string
	// HTTPS TLS certificate path
	HTTPSCertPath string
	// HTTPS TLS key path
	HTTPSKeyPath string
	// Client application relative path
	ClientAppURL string
	// Secret which is used for signing cookies
	Secret string
	// The entity which issued the certificate. Normally just domain name. See JWT docuymentation
	JWTIssuer string
	// JWT consumers. Normally just domain name. See JWT documentation
	JWTAudience string
	// Domain name where the cookie with token is valid
	CookieDomain string
	// Only SQL is supported
	StorageType storage.Type
	// If devmode is enabled, then CORS requests are allowed
	DevMode bool
}

// New creates new App config instance with pre-defined parameters
func New(ec EnvConfig) AppConfig {
	return AppConfig{
		ClientAppURL:  ec.ClientURL,
		CookieDomain:  ec.Domain,
		DSN:           ec.DSN,
		HTTPSCertPath: ec.HTTPSCertPath,
		HTTPSKeyPath:  ec.HTTPSKeyPath,
		JWTAudience:   ec.Domain,
		JWTIssuer:     ec.Domain,
		Secret:        ec.Secret,
		ServerAddr:    ec.ServerAddr,
		StorageType:   storage.TypeSQL,
		DevMode:       ec.DevMode,
	}
}
