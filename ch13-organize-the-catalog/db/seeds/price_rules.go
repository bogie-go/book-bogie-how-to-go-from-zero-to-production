package seeds

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"depot/app/models"
	"depot/db/adapter"
)

// SeedPriceRules loads the launch promotions: 20% off everything for the
// launch weekend, and buy two keycaps, get one free, next month.
func SeedPriceRules(ctx context.Context, s *adapter.Store) error {
	keycap, err := seededProduct(ctx, s, "Bogie Keycap")
	if err != nil {
		return err
	}
	now := time.Now()
	rules := []models.PriceRule{
		{
			Name: "Launch weekend: 20% off everything", Kind: "percent_off",
			Properties: mustJSON(models.PercentOff{Percent: 20}),
			Priority:   1, StartsAt: now, EndsAt: now.AddDate(0, 0, 3),
		},
		{
			Name: "Keycaps: buy 2, get 1 free", Kind: "buy_x_get_y",
			Properties: mustJSON(models.BuyXGetY{ProductID: keycap.ID, Buy: 2, Free: 1}),
			Priority:   2, StartsAt: now.AddDate(0, 1, 0), EndsAt: now.AddDate(0, 2, 0),
		},
		{
			Name: "$5 off orders of $40 or more", Kind: "amount_off",
			Properties: mustJSON(models.AmountOff{AmountCents: 500, MinSubtotalCents: 4000}),
			Priority:   3, StartsAt: now, EndsAt: now.AddDate(0, 1, 0),
		},
	}
	for _, r := range rules {
		if _, err := s.FindOrCreatePriceRuleByName(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// seededProduct finds a product from the products seed by its name.
func seededProduct(ctx context.Context, s *adapter.Store, name string) (models.Product, error) {
	for _, p := range products {
		if p.Name == name {
			return s.FindOrCreateProductByName(ctx, p)
		}
	}
	return models.Product{}, fmt.Errorf("seeds: no product named %q", name)
}

// mustJSON encodes a rule's settings. These structs always encode, so an
// error here is a bug in the seed itself.
func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
