package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"depot/app/models"
)

// ListActivePriceRules loads the rules running now, in the order they apply,
// each as the kind of rule its row names. A row that won't load is skipped and
// reported in skipped, so one bad row can't stop every cart.
func (s *Store) ListActivePriceRules(ctx context.Context) (
	rules []models.Rule, skipped []error, err error,
) {
	rows, err := s.q.ListActivePriceRules(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("store: active price rules: %w", err)
	}
	for _, row := range rows {
		rule, err := models.NewRule(priceRuleFromRow(row))
		if err != nil {
			skipped = append(skipped, fmt.Errorf("price rule %q: %w", row.Name, err))
			continue
		}
		rules = append(rules, rule)
	}
	return rules, skipped, nil
}

// FindOrCreatePriceRuleByName returns the rule named r.Name, creating it from
// r if there is none yet.
func (s *Store) FindOrCreatePriceRuleByName(ctx context.Context, r models.PriceRule) (models.PriceRule, error) {
	row, err := s.q.GetPriceRuleByName(ctx, r.Name)
	switch {
	case err == nil:
		return priceRuleFromRow(row), nil
	case errors.Is(err, pgx.ErrNoRows):
		return s.CreatePriceRule(ctx, r)
	default:
		return models.PriceRule{}, fmt.Errorf("store: find price rule %q: %w", r.Name, err)
	}
}
