package models

import (
	"fmt"
	"strings"
	"time"
)

// User is a user. Imports nothing third-party: ids are int64s here and
// become bigint at the store boundary.
type User struct {
	ID             int64
	Email          string
	PasswordDigest string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Validate is the model's validations. Wrap every
// failure in ErrInvalid so a controller can answer 422 without knowing the
// rule:
//
//	return fmt.Errorf("%w: title is required", ErrInvalid)
func (u User) Validate() error {
	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf("%w: email must be an email address", ErrInvalid)
	}
	if u.PasswordDigest == "" {
		return fmt.Errorf("%w: password is required", ErrInvalid)
	}
	return nil
}

// NormalizeEmail is how emails are compared: without surrounding spaces, and
// in lower case, so "Ada@Example.com" and "ada@example.com" are one account.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
