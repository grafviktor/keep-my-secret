// Package auth is used for generating access and refresh tokens,
// also it has methods to verify tokens validity
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

// Auth - struct contains all necessary information for JWT token generation
type Auth struct {
	// Issuer of the token. See JWT.io for more information
	Issuer string
	// Audience of the token. See JWT.io for more information
	Audience string
	// Secret of the token. Self-explanatory.
	Secret string
	// TokenExpiry is duration of the access token
	TokenExpiry time.Duration
	// RefreshExpiry is duration of the refresh token
	RefreshExpiry time.Duration
	// CookieDomain is domain of the cookie. Self-explanatory.
	CookieDomain string
	// CookiePath is path of the cookie. Self-explanatory. In our case it's always project root.
	CookiePath string
	//  CookieName is name of the cookie. Self-explanatory.
	CookieName string
}

// JWTUser - struct for storing user details. Only ID for the moment
type JWTUser struct {
	// ID contains user login which was used during registration process
	ID string `json:"id"`
}

// TokenPair - struct is used for marshaling tokens
type TokenPair struct {
	// AccessToken is a short-lived token which is used for accessing application resources.
	AccessToken string `json:"access_token"`
	// RefreshToken is a long-lived token which is used for querying access tokens
	RefreshToken string `json:"refresh_token"`
}

// Claims - is utilizing default jwt claims. A subject for further extension.
type Claims struct {
	jwt.RegisteredClaims
}

// New - creates new Auth struct and with application defined values
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

// GenerateTokenPair - create new Refresh and Access tokens
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

// GetRefreshCookie - create new refresh cookie which contains refresh token. Used when user logs in (or register).
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

// GetExpiredRefreshCookie - same as GetRefreshCookie,  but with expired time.  Used when user logs out.
// There is way delete existing cookie from browser, thus replacing the existing one with expired.
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

// JWTVerifier - used for verifying user tokens
type JWTVerifier struct{}

// VerifyAuthHeader - extract users token from HTTP header and verifies it.
//
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
