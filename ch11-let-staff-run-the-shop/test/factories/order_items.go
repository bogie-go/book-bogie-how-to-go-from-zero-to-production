package factories

import (
	"context"

	"depot/app/models"
	"depot/db/adapter"
)

// NewOrderItem is a valid order_item that is not saved: every attribute has a default,
// and "with" overrides what the spec cares about.
//
// Written once by `bogie g model` and yours to edit: add an attribute's
// default here when you add the column. Running the generator again leaves a
// changed file alone unless --force. The exception is `bogie model:sync`:
// when a migration adds or drops a column it merges the new default into this
// file and keeps your edits (a clash leaves conflict markers).
//
// Left at zero: OrderID, ProductID. The parent has to exist first, so set it from the
// parent's own factory:
//
//	factories.OrderItem(t, s, func(v *models.OrderItem) { v.OrderID = parent.ID })
func NewOrderItem(with ...func(*models.OrderItem)) models.OrderItem {
	record := models.OrderItem{
		Name:           Sequence("name "),
		Quantity:       1,
		UnitPriceCents: 1,
	}
	for _, f := range with {
		f(&record)
	}
	return record
}

// OrderItem saves a NewOrderItem through the real store and fails the spec if it
// cannot, so a spec reads as its arrangement and nothing else.
func OrderItem(t T, s *adapter.Store, with ...func(*models.OrderItem)) models.OrderItem {
	t.Helper()
	created, err := s.CreateOrderItem(context.Background(), NewOrderItem(with...))
	if err != nil {
		t.Fatal(err)
	}
	return created
}
