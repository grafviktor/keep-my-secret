package config

import "github.com/grafviktor/keep-my-secret/internal/storage"

type EnvConfig struct {
	Secret     string `env:"APP_SECRET" envDefault:"romeo romeo whiskey"`
	ServerAddr string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	// PSQL connection string (DSN)
	DSN string `env:"DSN" envDefault:"./kms.db"`
	// HTTPS TLS certificate path
	HTTPSCertPath string `env:"TLS_CERT_PATH" envDefault:"./tls/cert.pem"`
	// HTTPS TLS key path
	HTTPSKeyPath string `env:"TLS_KEY_PATH" envDefault:"./tls/key.pem"`
	Domain       string `env:"DOMAIN"     envDefault:"localhost"`
	ClientURL    string `env:"CLIENT_URL"     envDefault:"/"`
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
	Secret       string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	StorageType  storage.Type
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
	}
}
