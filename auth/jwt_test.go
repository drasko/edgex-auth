package auth_test

import (
	"testing"

	"github.com/drasko/go-auth/domain"
)

func TestSubjectOf(t *testing.T) {
	user := "test-user"
	key, _ := domain.CreateKey(user)

	cases := []struct {
		user   string
		token  string
		hasErr bool
	}{
		{user, key, false},
		{user, "", true},
	}

	for i, c := range cases {
		subject, err := domain.SubjectOf(c.token)
		if c.hasErr && err == nil {
			t.Errorf("case %d: expected error to be thrown", i+1)
			continue
		}

		if !c.hasErr && err != nil {
			t.Errorf("case %d: didn't expect an error to be thrown", i+1)
			continue
		}

		if err == nil && c.user != subject {
			t.Errorf("case %d: expected user %s got %s", i+1, c.user, subject)
		}
	}
}
