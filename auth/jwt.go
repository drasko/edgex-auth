package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const Issuer string = "mainflux"

var secretKey string = "mainflux-api-key"

// SetSecretKey sets the secret key that will be used for decoding and encoding
// tokens. If not invoked, a default key will be used instead.
func SetSecretKey(key string) {
	secretKey = key
}

// DecodeJWT decodes jwt token
func DecodeJwt(key string) (*jwt.StandardClaims, error) {
	c := jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(
		key,
		&c,
		func(val *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Validate the token and return the custom claims
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		err := errors.New("Decoded invalid token")
		return nil, err
	}

	return claims, nil
}

// CreateKey creates a JSON Web Token with a given subject.
func CreateKey(subject string) (string, error) {
	claims := jwt.StandardClaims{
		Issuer:   Issuer,
		IssuedAt: time.Now().UTC().Unix(),
		Subject:  subject,
	}

	key := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := key.SignedString([]byte(secretKey))
	if err != nil {
		return "", &AuthError{http.StatusInternalServerError, err.Error()}
	}

	return raw, nil
}
