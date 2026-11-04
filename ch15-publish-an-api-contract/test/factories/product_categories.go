package factories

import (
	"context"

	"depot/app/models"
	"depot/db/adapter"
)

// NewProductCategory is a valid product_category that is not saved: every attribute has a default,
// and "with" overrides what the spec cares about.
//
// Written once by `bogie g model` and yours to edit: add an attribute's
// default here when you add the column. Running the generator again leaves a
// changed file alone unless --force. The exception is `bogie model:sync`:
// when a migration adds or drops a column it merges the new default into this
// file and keeps your edits (a clash leaves conflict markers).
//
// Left at zero: ProductID, CategoryID. The parent has to exist first, so set it from the
// parent's own factory:
//
//	factories.ProductCategory(t, s, func(v *models.ProductCategory) { v.ProductID = parent.ID })
func NewProductCategory(with ...func(*models.ProductCategory)) models.ProductCategory {
	record := models.ProductCategory{}
	for _, f := range with {
		f(&record)
	}
	return record
}

// ProductCategory saves a NewProductCategory through the real store and fails the spec if it
// cannot, so a spec reads as its arrangement and nothing else.
func ProductCategory(t T, s *adapter.Store, with ...func(*models.ProductCategory)) models.ProductCategory {
	t.Helper()
	created, err := s.CreateProductCategory(context.Background(), NewProductCategory(with...))
	if err != nil {
		t.Fatal(err)
	}
	return created
}
