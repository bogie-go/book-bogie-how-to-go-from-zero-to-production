package factories

import (
	"context"

	"depot/app/models"
	"depot/db/adapter"
)

// NewCartItem is a valid cart_item that is not saved: every attribute has a default,
// and "with" overrides what the spec cares about.
//
// Written once by `bogie g model` and yours to edit: add an attribute's
// default here when you add the column. Running the generator again leaves a
// changed file alone unless --force. The exception is `bogie model:sync`:
// when a migration adds or drops a column it merges the new default into this
// file and keeps your edits (a clash leaves conflict markers).
//
// Left at zero: CartID, ProductID. The parent has to exist first, so set it from the
// parent's own factory:
//
//	factories.CartItem(t, s, func(v *models.CartItem) { v.CartID = parent.ID })
func NewCartItem(with ...func(*models.CartItem)) models.CartItem {
	record := models.CartItem{
		Quantity: 1,
	}
	for _, f := range with {
		f(&record)
	}
	return record
}

// CartItem saves a NewCartItem through the real store and fails the spec if it
// cannot, so a spec reads as its arrangement and nothing else.
func CartItem(t T, s *adapter.Store, with ...func(*models.CartItem)) models.CartItem {
	t.Helper()
	created, err := s.CreateCartItem(context.Background(), NewCartItem(with...))
	if err != nil {
		t.Fatal(err)
	}
	return created
}
