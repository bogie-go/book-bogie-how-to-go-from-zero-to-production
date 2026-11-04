package factories

import (
	"context"

	"depot/app/models"
	"depot/db/adapter"
)

// NewProduct is a valid product that is not saved: every attribute has a default,
// and "with" overrides what the spec cares about.
//
// Written once by `bogie g model` and yours to edit: add an attribute's
// default here when you add the column. Running the generator again leaves a
// changed file alone unless --force. The exception is `bogie model:sync`:
// when a migration adds or drops a column it merges the new default into this
// file and keeps your edits (a clash leaves conflict markers).
func NewProduct(with ...func(*models.Product)) models.Product {
	record := models.Product{
		Name:        Sequence("name "),
		Description: Sequence("description "),
		PriceCents:  1,
		Stock:       1,
		ImageUrl:    Sequence("image url "),
	}
	for _, f := range with {
		f(&record)
	}
	return record
}

// Product saves a NewProduct through the real store and fails the spec if it
// cannot, so a spec reads as its arrangement and nothing else.
func Product(t T, s *adapter.Store, with ...func(*models.Product)) models.Product {
	t.Helper()
	created, err := s.CreateProduct(context.Background(), NewProduct(with...))
	if err != nil {
		t.Fatal(err)
	}
	return created
}
