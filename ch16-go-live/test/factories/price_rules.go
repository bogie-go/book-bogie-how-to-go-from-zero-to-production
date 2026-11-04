package factories

import (
	"context"
	"time"

	"depot/app/models"
	"depot/db/adapter"
)

// NewPriceRule is a valid price_rule that is not saved: every attribute has a default,
// and "with" overrides what the spec cares about.
//
// Written once by `bogie g model` and yours to edit: add an attribute's
// default here when you add the column. Running the generator again leaves a
// changed file alone unless --force. The exception is `bogie model:sync`:
// when a migration adds or drops a column it merges the new default into this
// file and keeps your edits (a clash leaves conflict markers).
//
// Placeholders: the datetime defaults (StartsAt, EndsAt) are one fixed instant. Set real
// times where a spec depends on them, such as a window that starts before it ends.
func NewPriceRule(with ...func(*models.PriceRule)) models.PriceRule {
	record := models.PriceRule{
		Name:       Sequence("name "),
		Kind:       "percent_off",
		Properties: []byte(`{"percent":10}`),
		Priority:   1,
		StartsAt:   time.Now().Add(-time.Hour),
		EndsAt:     time.Now().Add(24 * time.Hour),
	}
	for _, f := range with {
		f(&record)
	}
	return record
}

// PriceRule saves a NewPriceRule through the real store and fails the spec if it
// cannot, so a spec reads as its arrangement and nothing else.
func PriceRule(t T, s *adapter.Store, with ...func(*models.PriceRule)) models.PriceRule {
	t.Helper()
	created, err := s.CreatePriceRule(context.Background(), NewPriceRule(with...))
	if err != nil {
		t.Fatal(err)
	}
	return created
}
