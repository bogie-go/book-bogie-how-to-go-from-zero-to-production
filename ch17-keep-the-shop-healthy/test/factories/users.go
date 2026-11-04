package factories

import (
	"context"

	"depot/app/models"
	"depot/db/adapter"
)

// NewUser is a valid user that is not saved: every attribute has a default,
// and "with" overrides what the spec cares about.
//
// Written once by `bogie g model` and yours to edit: add an attribute's
// default here when you add the column. Running the generator again leaves a
// changed file alone unless --force. The exception is `bogie model:sync`:
// when a migration adds or drops a column it merges the new default into this
// file and keeps your edits (a clash leaves conflict markers).
func NewUser(with ...func(*models.User)) models.User {
	record := models.User{
		Email:          Sequence("customer-") + "@example.com",
		PasswordDigest: Sequence("password digest "),
	}
	for _, f := range with {
		f(&record)
	}
	return record
}

// User saves a NewUser through the real store and fails the spec if it
// cannot, so a spec reads as its arrangement and nothing else.
func User(t T, s *adapter.Store, with ...func(*models.User)) models.User {
	t.Helper()
	created, err := s.CreateUser(context.Background(), NewUser(with...))
	if err != nil {
		t.Fatal(err)
	}
	return created
}
