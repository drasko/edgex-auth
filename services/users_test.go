package services_test

import (
	"net/http"
	"testing"

	"github.com/drasko/go-auth/domain"
	"github.com/drasko/go-auth/services"
)

func TestRegisterUser(t *testing.T) {
	var (
		username = "register-username"
		password = "register-password"
	)

	cases := []struct {
		username string
		password string
		code     int
	}{
		{username, password, http.StatusOK},
		{username, password, http.StatusConflict},
		{username, "", http.StatusBadRequest},
		{"", password, http.StatusBadRequest},
	}

	for i, c := range cases {
		_, err := services.RegisterUser(c.username, c.password)
		if err != nil {
			auth := err.(*domain.AuthError)
			if auth.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, auth.Code)
			}
		}
	}
}

func TestLogin(t *testing.T) {
	var (
		username = "login-username"
		password = "login-password"
	)

	services.RegisterUser(username, password)

	cases := []struct {
		username string
		password string
		code     int
	}{
		{username, password, http.StatusOK},
		{username, "", http.StatusBadRequest},
		{"", password, http.StatusBadRequest},
		{"bad", password, http.StatusForbidden},
		{username, "bad", http.StatusForbidden},
	}

	for i, c := range cases {
		_, err := services.Login(c.username, c.password)
		if err != nil {
			auth := err.(*domain.AuthError)
			if auth.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, auth.Code)
			}
		}
	}
}
