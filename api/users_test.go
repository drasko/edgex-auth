package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	cases := []struct {
		body string
		code int
	}{
		{`{"username":"test","password":"test"}`, http.StatusCreated},
		{"malformed body", http.StatusBadRequest},
		{`{"username":"","password":"test"}`, http.StatusBadRequest},
		{`{"username":"test","password":""}`, http.StatusBadRequest},
		{`{"username":"test","password":"test"}`, http.StatusConflict},
	}

	url := fmt.Sprintf("%s/users", ts.URL)

	for i, c := range cases {
		b := strings.NewReader(c.body)

		res, err := http.Post(url, "application/json", b)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}
	}
}

func TestLoginUser(t *testing.T) {
	cases := []struct {
		body string
		code int
	}{
		{`{"username":"test","password":"test"}`, http.StatusCreated},
		{"malformed body", http.StatusBadRequest},
		{`{"username":"","password":""}`, http.StatusBadRequest},
		{`{"username":"","password":"test"}`, http.StatusBadRequest},
		{`{"username":"test","password":""}`, http.StatusBadRequest},
		{`{"username":"bad","password":"test"}`, http.StatusForbidden},
		{`{"username":"test","password":"bad"}`, http.StatusForbidden},
	}

	url := fmt.Sprintf("%s/sessions", ts.URL)

	for i, c := range cases {
		b := strings.NewReader(c.body)

		res, err := http.Post(url, "application/json", b)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}

		defer res.Body.Close()
	}
}
