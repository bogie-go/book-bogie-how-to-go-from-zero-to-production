package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RoleStaff is the role of the shop's staff. A customer's role is empty.
const RoleStaff = "staff"

// User is a user. Imports nothing third-party: ids are int64s here and
// become bigint at the store boundary.
type User struct {
	ID             int64
	Email          string
	PasswordDigest string
	Role           string
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
	if u.Role != "" && u.Role != RoleStaff {
		return fmt.Errorf("%w: role must be empty or %q", ErrInvalid, RoleStaff)
	}
	return nil
}

// NormalizeEmail is how emails are compared: without surrounding spaces, and
// in lower case, so "Ada@Example.com" and "ada@example.com" are one account.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// UserSubject is the token subject for a customer's account.
func UserSubject(id int64) string {
	return "user:" + strconv.FormatInt(id, 10)
}

// UserIDFromSubject is the customer a token's subject names, if it names one.
func UserIDFromSubject(subject string) (int64, bool) {
	rest, ok := strings.CutPrefix(subject, "user:")
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseInt(rest, 10, 64)
	return id, err == nil
}
