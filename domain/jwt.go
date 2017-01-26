package domain

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const issuer string = "mainflux"

var secretKey string = "mainflux-api-key"

// SetSecretKey sets the secret key that will be used for decoding and encoding
// tokens. If not invoked, a default key will be used instead.
func SetSecretKey(key string) {
	secretKey = key
}

// SubjectOf extracts token's subject.
func SubjectOf(key string) (string, error) {
	claims := jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(
		key,
		&claims,
		func(val *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil || !token.Valid || claims.Issuer != issuer {
		return "", &AuthError{http.StatusForbidden, err.Error()}
	}

	return claims.Subject, nil
}

// CreateKey creates a JSON Web Token with a given subject.
func CreateKey(subject string) (string, error) {
	claims := jwt.StandardClaims{
		Issuer:   issuer,
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
