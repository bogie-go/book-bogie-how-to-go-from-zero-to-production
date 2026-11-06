package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"depot/app/models"
	"depot/db/sqlc"
)

// PriceRules, as the models see them. Each method converts models.PriceRule to and from
// the row type sqlc generated, and maps pgx's "no rows" to models.ErrNotFound
// so callers never import pgx to ask a question.
//
// Written by `bogie g model`: overwritten whole on the next run. A
// hand-written addition goes in price_rules_extra.go, not here.

// CreatePriceRule validates a price_rule, inserts it, and returns it with its id and
// timestamps. Validation runs on save, so every caller
// gets it and none can forget it.
func (s *Store) CreatePriceRule(ctx context.Context, p models.PriceRule) (models.PriceRule, error) {
	if err := p.Validate(); err != nil {
		return models.PriceRule{}, err
	}
	row, err := s.q.CreatePriceRule(ctx, sqlc.CreatePriceRuleParams{Name: p.Name, Kind: p.Kind, Properties: p.Properties, Priority: p.Priority, StartsAt: pgtype.Timestamptz{Time: p.StartsAt, Valid: true}, EndsAt: pgtype.Timestamptz{Time: p.EndsAt, Valid: true}})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.PriceRule{}, fmt.Errorf("%w: name is already taken", models.ErrInvalid)
		}
		return models.PriceRule{}, fmt.Errorf("store: create price_rule: %w", err)
	}
	return priceRuleFromRow(row), nil
}

// GetPriceRule finds one price_rule. An id that is not a valid id is simply not found.
func (s *Store) GetPriceRule(ctx context.Context, id string) (models.PriceRule, error) {
	pk, err := parseID(id)
	if err != nil {
		return models.PriceRule{}, models.ErrNotFound
	}
	row, err := s.q.GetPriceRule(ctx, pk)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.PriceRule{}, models.ErrNotFound
	case err != nil:
		return models.PriceRule{}, fmt.Errorf("store: get price_rule %s: %w", id, err)
	}
	return priceRuleFromRow(row), nil
}

// ListPriceRules returns the newest price_rules first, at most limit of them.
func (s *Store) ListPriceRules(ctx context.Context, limit int32) ([]models.PriceRule, error) {
	rows, err := s.q.ListPriceRules(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list price_rules: %w", err)
	}
	out := make([]models.PriceRule, 0, len(rows))
	for _, r := range rows {
		out = append(out, priceRuleFromRow(r))
	}
	return out, nil
}

// UpdatePriceRule validates a price_rule, writes its columns, and returns the stored row.
func (s *Store) UpdatePriceRule(ctx context.Context, p models.PriceRule) (models.PriceRule, error) {
	if err := p.Validate(); err != nil {
		return models.PriceRule{}, err
	}
	pk := p.ID
	row, err := s.q.UpdatePriceRule(ctx, sqlc.UpdatePriceRuleParams{ID: pk, Name: p.Name, Kind: p.Kind, Properties: p.Properties, Priority: p.Priority, StartsAt: pgtype.Timestamptz{Time: p.StartsAt, Valid: true}, EndsAt: pgtype.Timestamptz{Time: p.EndsAt, Valid: true}})
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.PriceRule{}, models.ErrNotFound
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		return models.PriceRule{}, fmt.Errorf("%w: name is already taken", models.ErrInvalid)
	case err != nil:
		return models.PriceRule{}, fmt.Errorf("store: update price_rule %v: %w", p.ID, err)
	}
	return priceRuleFromRow(row), nil
}

// DeletePriceRule removes a price_rule, and says so if there was none.
func (s *Store) DeletePriceRule(ctx context.Context, id string) error {
	pk, err := parseID(id)
	if err != nil {
		return models.ErrNotFound
	}
	n, err := s.q.DeletePriceRule(ctx, pk)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("%w: price_rule is still in use", models.ErrInvalid)
		}
		return fmt.Errorf("store: delete price_rule %s: %w", id, err)
	}
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func priceRuleFromRow(r sqlc.PriceRule) models.PriceRule {
	return models.PriceRule{
		ID:         r.ID,
		Name:       r.Name,
		Kind:       r.Kind,
		Properties: r.Properties,
		Priority:   r.Priority,
		StartsAt:   r.StartsAt.Time,
		EndsAt:     r.EndsAt.Time,
		CreatedAt:  r.CreatedAt.Time,
		UpdatedAt:  r.UpdatedAt.Time,
	}
}
