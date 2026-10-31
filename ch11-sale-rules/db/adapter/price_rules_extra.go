package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"depot/app/models"
)

// ListActivePriceRules lists the rules running now, in the order they apply.
func (s *Store) ListActivePriceRules(ctx context.Context) ([]models.PriceRule, error) {
	rows, err := s.q.ListActivePriceRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: active price rules: %w", err)
	}
	rules := make([]models.PriceRule, 0, len(rows))
	for _, r := range rows {
		rules = append(rules, priceRuleFromRow(r))
	}
	return rules, nil
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
