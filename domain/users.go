package domain

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents one Mainflux user account.
type User struct {
	Id        string `json:"id"`
	Username  string `json:"-"`
	Password  string `json:"-"`
	MasterKey string `json:"key"`
}

// CreateUser creates new user account based on provided username and password.
// The account is assigned with one master key - a key with all permissions on
// all owned resources regardless of their type. Provided password in encrypted
// using bcrypt algorithm.
func CreateUser(username, password string) (User, error) {
	user := User{
		Id:       uuid.NewV4().String(),
		Username: username,
	}

	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, &AuthError{Code: http.StatusInternalServerError}
	}

	user.Password = string(p)
	user.MasterKey, err = CreateKey(user.Id)
	if err != nil {
		return user, err
	}

	return user, nil
}

// CheckPassword tries to determine whether or not the submitted password
// matches the one stored (and hashed) during registration. An error will be
// used to indicate an invalid password.
func CheckPassword(plain, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
