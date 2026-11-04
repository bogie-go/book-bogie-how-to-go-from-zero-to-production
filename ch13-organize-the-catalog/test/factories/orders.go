package factories

import (
	"context"

	"depot/app/models"
	"depot/db/adapter"
)

// NewOrder is a valid order that is not saved: every attribute has a default,
// and "with" overrides what the spec cares about.
//
// Written once by `bogie g model` and yours to edit: add an attribute's
// default here when you add the column. Running the generator again leaves a
// changed file alone unless --force. The exception is `bogie model:sync`:
// when a migration adds or drops a column it merges the new default into this
// file and keeps your edits (a clash leaves conflict markers).
//
// Left at zero: UserID. The parent has to exist first, so set it from the
// parent's own factory:
//
//	factories.Order(t, s, func(v *models.Order) { v.UserID = parent.ID })
func NewOrder(with ...func(*models.Order)) models.Order {
	record := models.Order{
		Status:        models.OrderPlaced,
		SubtotalCents: 1000,
		DiscountCents: 0,
		TotalCents:    1000,
	}
	for _, f := range with {
		f(&record)
	}
	return record
}

// Order saves a NewOrder through the real store and fails the spec if it
// cannot, so a spec reads as its arrangement and nothing else.
func Order(t T, s *adapter.Store, with ...func(*models.Order)) models.Order {
	t.Helper()
	created, err := s.CreateOrder(context.Background(), NewOrder(with...))
	if err != nil {
		t.Fatal(err)
	}
	return created
}
