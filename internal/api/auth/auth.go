package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grafviktor/keep-my-secret/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// CookieName which start from '__Host-' will NOT be set if domain specified,
	// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies#cookie_prefixes
	// Also, a cookie with such prefix, cannot be set over http connection
	CookieName   = "__Host-refresh_token"
	cookieSecure = true

	// SameSiteStrictMode will not allow to set cookie for CORS (Cross-origin resource sharing) connections.
	// siteMode     = http.SameSiteStrictMode
	// Bypassing set cookie request in CORS connections. However, you also must be sure that cookie is "secure: true"
	siteMode = http.SameSiteNoneMode
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type JWTUser struct {
	ID string `json:"id"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
}

func New(ac config.AppConfig) Auth {
	return Auth{
		Issuer:        ac.JWTIssuer,
		Audience:      ac.JWTAudience,
		Secret:        ac.Secret,
		TokenExpiry:   time.Minute * 10,
		RefreshExpiry: time.Hour * 24,
		CookieDomain:  ac.CookieDomain,
		CookiePath:    "/",
		CookieName:    CookieName,
	}
}

func (auth Auth) GenerateTokenPair(user *JWTUser) (TokenPair, error) { // pair for token and refresh token
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a accessToken and set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID       // id of the user in a database
	claims["aud"] = auth.Audience // audience
	claims["iss"] = auth.Issuer
	claims["iat"] = time.Now().UTC().Unix()                       // issued at
	claims["typ"] = "JWT"                                         // type
	claims["exp"] = time.Now().UTC().Add(auth.TokenExpiry).Unix() // expiry

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(auth.Secret))
	if err != nil {
		return TokenPair{}, err
	}

	// Create a refreshToken and set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = user.ID                                         // id of the user in a database
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()                         // issued at
	refreshTokenClaims["exp"] = time.Now().UTC().Add(auth.RefreshExpiry).Unix() // expiry

	// Create signed refresh token
	signedRefreshToken, err := token.SignedString([]byte(auth.Secret))
	if err != nil {
		return TokenPair{}, err
	}

	// Create token pairs and populate with signed tokens
	tokenPairs := TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Return TokenPair
	return tokenPairs, nil
}

func (auth Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     auth.CookieName,
		Path:     auth.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().Add(auth.RefreshExpiry),
		MaxAge:   int(auth.RefreshExpiry.Seconds()),
		SameSite: siteMode,
		HttpOnly: true,
		Secure:   cookieSecure,
	}
}

// GetExpiredRefreshCookie - For logging out
func (auth Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     auth.CookieName,
		Path:     auth.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: siteMode,
		HttpOnly: true,
		Secure:   cookieSecure,
	}
}

type JWTVerifier struct{}

//nolint:lll
func (t JWTVerifier) VerifyAuthHeader(config config.AppConfig, w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	headerParts := strings.Split(authHeader, " ") // "Bearer Token"
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", nil, errors.New("invalid auth header")
	}

	token := headerParts[1]
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return []byte(config.Secret), nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}

		return "", nil, err
	}

	if claims.Issuer != config.JWTIssuer {
		return "", nil, errors.New("invalid issuer")
	}

	return token, claims, nil
}
