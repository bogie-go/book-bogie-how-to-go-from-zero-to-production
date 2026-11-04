// Package password hashes and checks customers' passwords, so the rest of
// the app never touches bcrypt directly.
package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"depot/app/models"
)

const (
	minLength = 8
	maxLength = 72 // bcrypt only looks at the first 72 bytes
)

// Hash turns a password into a digest that can be stored, refusing a
// password that's too short or too long.
func Hash(plain string) (string, error) {
	if len(plain) < minLength || len(plain) > maxLength {
		return "", fmt.Errorf("%w: password must be %d to %d bytes long", models.ErrInvalid, minLength, maxLength)
	}
	digest, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password: hash: %w", err)
	}
	return string(digest), nil
}

// Check reports whether plain is the password that digest was made from.
func Check(digest, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(digest), []byte(plain)) == nil
}

// dummyDigest is a digest of nothing anyone will type, for DummyCheck.
var dummyDigest = mustHash("not anyone's password")

// mustHash is Hash for a password known to be valid; an error is a bug.
func mustHash(plain string) string {
	digest, err := Hash(plain)
	if err != nil {
		panic(err)
	}
	return digest
}

// DummyCheck does the work of a Check and throws the answer away. Sign-in
// calls it for an email with no account, so that answer takes as long as a
// wrong password does.
func DummyCheck(plain string) {
	Check(dummyDigest, plain)
}
