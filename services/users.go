package services

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/domain"
)

// RegisterUser invokes creation of new user account based on provided username
// and password.
func RegisterUser(username, password string) (domain.User, error) {
	var user domain.User

	if username == "" || password == "" {
		return user, &domain.AuthError{Code: http.StatusBadRequest}
	}

	c := cache.Get()
	defer c.Close()

	userKey := fmt.Sprintf("auth:%s:%s:profile", domain.UserType, username)
	if exists, _ := redis.Bool(c.Do("EXISTS", userKey)); exists {
		return user, &domain.AuthError{Code: http.StatusConflict}
	}

	user, err := domain.CreateUser(username, password)
	if err != nil {
		return user, err
	}

	masterKey := fmt.Sprintf("auth:%s:%s:master", domain.UserType, user.Id)

	c.Send("WATCH", userKey, masterKey)
	c.Send("MULTI")
	c.Send("HMSET", userKey, "password", user.Password, "id", user.Id)
	c.Send("SET", masterKey, user.MasterKey)
	_, err = c.Do("EXEC")
	if err != nil {
		return user, &domain.AuthError{Code: http.StatusInternalServerError}
	}

	return user, nil
}

// Login retrieves user's master key when invoked with valid username and
// password.
func Login(username, password string) (domain.User, error) {
	var user domain.User

	if username == "" || password == "" {
		return user, &domain.AuthError{Code: http.StatusBadRequest}
	}

	c := cache.Get()
	defer c.Close()

	cKey := fmt.Sprintf("auth:%s:%s:profile", domain.UserType, username)

	items, err := redis.Strings(c.Do("HMGET", cKey, "id", "password"))
	if err != nil {
		return user, &domain.AuthError{Code: http.StatusForbidden}
	}

	if err := domain.CheckPassword(password, items[1]); err != nil {
		return user, &domain.AuthError{Code: http.StatusForbidden}
	}

	user.Id = items[0]
	cKey = fmt.Sprintf("auth:%s:%s:master", domain.UserType, user.Id)

	user.MasterKey, _ = redis.String(c.Do("GET", cKey))
	if user.MasterKey == "" {
		return user, &domain.AuthError{Code: http.StatusForbidden}
	}

	return user, nil
}
