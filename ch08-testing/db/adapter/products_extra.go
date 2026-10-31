package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"depot/app/models"
)

// FindOrCreateProductByName returns the product named p.Name, creating it
// from p if there is none yet.
func (s *Store) FindOrCreateProductByName(ctx context.Context, p models.Product) (models.Product, error) {
	row, err := s.q.GetProductByName(ctx, p.Name)
	switch {
	case err == nil:
		return productFromRow(row), nil
	case errors.Is(err, pgx.ErrNoRows):
		return s.CreateProduct(ctx, p)
	default:
		return models.Product{}, fmt.Errorf("store: find product %q: %w", p.Name, err)
	}
}
