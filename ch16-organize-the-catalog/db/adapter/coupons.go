package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"depot/app/models"
	"depot/db/sqlc"
)

// Coupons, as the models see them. Each method converts models.Coupon to and from
// the row type sqlc generated, and maps pgx's "no rows" to models.ErrNotFound
// so callers never import pgx to ask a question.
//
// Written by `bogie g model`: overwritten whole on the next run. A
// hand-written addition goes in coupons_extra.go, not here.

// CreateCoupon validates a coupon, inserts it, and returns it with its id and
// timestamps. Validation runs on save, so every caller
// gets it and none can forget it.
func (s *Store) CreateCoupon(ctx context.Context, c models.Coupon) (models.Coupon, error) {
	if err := c.Validate(); err != nil {
		return models.Coupon{}, err
	}
	row, err := s.q.CreateCoupon(ctx, sqlc.CreateCouponParams{Code: c.Code, PriceRuleID: c.PriceRuleID})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.Coupon{}, fmt.Errorf("%w: code is already taken", models.ErrInvalid)
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return models.Coupon{}, fmt.Errorf("%w: price_rule_id does not exist", models.ErrInvalid)
		}
		return models.Coupon{}, fmt.Errorf("store: create coupon: %w", err)
	}
	return couponFromRow(row), nil
}

// GetCoupon finds one coupon. An id that is not a valid id is simply not found.
func (s *Store) GetCoupon(ctx context.Context, id string) (models.Coupon, error) {
	pk, err := parseID(id)
	if err != nil {
		return models.Coupon{}, models.ErrNotFound
	}
	row, err := s.q.GetCoupon(ctx, pk)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Coupon{}, models.ErrNotFound
	case err != nil:
		return models.Coupon{}, fmt.Errorf("store: get coupon %s: %w", id, err)
	}
	return couponFromRow(row), nil
}

// ListCoupons returns the newest coupons first, at most limit of them.
func (s *Store) ListCoupons(ctx context.Context, limit int32) ([]models.Coupon, error) {
	rows, err := s.q.ListCoupons(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list coupons: %w", err)
	}
	out := make([]models.Coupon, 0, len(rows))
	for _, r := range rows {
		out = append(out, couponFromRow(r))
	}
	return out, nil
}

// UpdateCoupon validates a coupon, writes its columns, and returns the stored row.
func (s *Store) UpdateCoupon(ctx context.Context, c models.Coupon) (models.Coupon, error) {
	if err := c.Validate(); err != nil {
		return models.Coupon{}, err
	}
	pk := c.ID
	row, err := s.q.UpdateCoupon(ctx, sqlc.UpdateCouponParams{ID: pk, Code: c.Code, PriceRuleID: c.PriceRuleID})
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Coupon{}, models.ErrNotFound
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		return models.Coupon{}, fmt.Errorf("%w: code is already taken", models.ErrInvalid)
	case errors.As(err, &pgErr) && pgErr.Code == "23503":
		return models.Coupon{}, fmt.Errorf("%w: price_rule_id does not exist", models.ErrInvalid)
	case err != nil:
		return models.Coupon{}, fmt.Errorf("store: update coupon %v: %w", c.ID, err)
	}
	return couponFromRow(row), nil
}

// DeleteCoupon removes a coupon, and says so if there was none.
func (s *Store) DeleteCoupon(ctx context.Context, id string) error {
	pk, err := parseID(id)
	if err != nil {
		return models.ErrNotFound
	}
	n, err := s.q.DeleteCoupon(ctx, pk)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("%w: coupon is still in use", models.ErrInvalid)
		}
		return fmt.Errorf("store: delete coupon %s: %w", id, err)
	}
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func couponFromRow(r sqlc.Coupon) models.Coupon {
	return models.Coupon{
		ID:          r.ID,
		Code:        r.Code,
		PriceRuleID: r.PriceRuleID,
		CreatedAt:   r.CreatedAt.Time,
		UpdatedAt:   r.UpdatedAt.Time,
	}
}
