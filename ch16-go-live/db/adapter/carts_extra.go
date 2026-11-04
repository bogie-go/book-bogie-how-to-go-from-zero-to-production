package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"depot/app/models"
)

// FindOrCreateCartByToken returns the cart with this token, creating it if
// there is none. It's one statement, so concurrent callers get the same cart.
func (s *Store) FindOrCreateCartByToken(ctx context.Context, token string) (models.Cart, error) {
	row, err := s.q.UpsertCartByToken(ctx, token)
	if err != nil {
		return models.Cart{}, fmt.Errorf("store: find or create cart: %w", err)
	}
	return cartFromRow(row), nil
}

// MergeCart moves every line of the cart with token from into the cart with
// token to, then removes the old cart. It's best-effort: a line that can't
// move, say because it's out of stock, is skipped and reported, and when
// there's no old cart, there's simply nothing to merge.
func (s *Store) MergeCart(ctx context.Context, from, to string) ([]error, error) {
	if from == to {
		return nil, nil // a cart merged into itself would only be emptied
	}
	old, err := s.q.GetCartByToken(ctx, from)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("store: merge cart: %w", err)
	}
	target, err := s.FindOrCreateCartByToken(ctx, to)
	if err != nil {
		return nil, err
	}
	lines, err := s.CartContents(ctx, old.ID)
	if err != nil {
		return nil, err
	}
	var skipped []error
	for _, l := range lines {
		if err := s.AddCartItem(ctx, target.ID, l.ProductID, l.Quantity); err != nil {
			skipped = append(skipped, fmt.Errorf("product %d: %w", l.ProductID, err))
		}
	}
	if _, err := s.q.DeleteCartItemsForCart(ctx, old.ID); err != nil {
		return skipped, fmt.Errorf("store: merge cart: %w", err)
	}
	if _, err := s.q.DeleteCart(ctx, old.ID); err != nil {
		return skipped, fmt.Errorf("store: merge cart: %w", err)
	}
	return skipped, nil
}
