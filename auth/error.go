package auth

import "fmt"

// AuthError is used as a wrapper for all types of errors originated
// throughout the service. Code will be used as API response status.
type AuthError struct {
	Code    int
	Message string
}

func (e *AuthError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
