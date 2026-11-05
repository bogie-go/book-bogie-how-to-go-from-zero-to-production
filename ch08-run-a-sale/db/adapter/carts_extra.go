package adapter

import (
	"context"
	"fmt"

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
