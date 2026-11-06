package factories

import (
	"context"

	"depot/app/models"
	"depot/db/adapter"
)

// NewCoupon is a valid coupon that is not saved: every attribute has a default,
// and "with" overrides what the spec cares about.
//
// Written once by `bogie g model` and yours to edit: add an attribute's
// default here when you add the column. Running the generator again leaves a
// changed file alone unless --force. The exception is `bogie model:sync`:
// when a migration adds or drops a column it merges the new default into this
// file and keeps your edits (a clash leaves conflict markers).
//
// Left at zero: PriceRuleID. The parent has to exist first, so set it from the
// parent's own factory:
//
//	factories.Coupon(t, s, func(v *models.Coupon) { v.PriceRuleID = parent.ID })
func NewCoupon(with ...func(*models.Coupon)) models.Coupon {
	record := models.Coupon{
		Code: Sequence("CODE-"),
	}
	for _, f := range with {
		f(&record)
	}
	return record
}

// Coupon saves a NewCoupon through the real store and fails the spec if it
// cannot, so a spec reads as its arrangement and nothing else.
func Coupon(t T, s *adapter.Store, with ...func(*models.Coupon)) models.Coupon {
	t.Helper()
	created, err := s.CreateCoupon(context.Background(), NewCoupon(with...))
	if err != nil {
		t.Fatal(err)
	}
	return created
}
