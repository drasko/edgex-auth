package domain_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/drasko/go-auth/domain"
)

func TestCreateUser(t *testing.T) {
	cases := []struct {
		username string
		password string
	}{
		{"x", "x"},
		{"y", "y"},
	}

	for i, c := range cases {
		user, err := domain.CreateUser(c.username, c.username)
		if err != nil {
			_, ok := err.(*domain.AuthError)
			if !ok {
				t.Errorf("case %d: all errors must be AuthError", i+1)
			}
		}

		if user.Username != c.username {
			t.Errorf("case %d: expected %s got %s", i+1, c.username, user.Username)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.password))
		if err != nil {
			t.Errorf("case %d: invalid password", i+1)
		}

		subject, err := domain.SubjectOf(user.MasterKey)
		if err != nil {
			t.Errorf("case %d: invalid master key", i+1)
		}

		if user.Id != subject {
			t.Errorf("case %d: expected %s got %s", i+1, subject, user.Id)
		}
	}
}

func TestCheckPassword(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

	cases := []struct {
		plain  string
		hashed string
		hasErr bool
	}{
		{"test", string(hashed), false},
		{"bad", string(hashed), true},
	}

	for i, c := range cases {
		err := domain.CheckPassword(c.plain, c.hashed)

		hasErr := err != nil
		if c.hasErr != hasErr {
			t.Errorf("case %d: expected %t got %t", i, c.hasErr, hasErr)
		}
	}
}
Contact GitHub API Training Shop Blog About

