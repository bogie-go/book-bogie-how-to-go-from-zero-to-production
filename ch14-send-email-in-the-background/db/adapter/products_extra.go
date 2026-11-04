package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"depot/app/models"
	"depot/db/sqlc"
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

// ArchiveProduct takes a product out of the shop. It touches only
// archived_at, so it can't undo a checkout's change to stock.
func (s *Store) ArchiveProduct(ctx context.Context, id string) (models.Product, error) {
	pk, err := parseID(id)
	if err != nil {
		return models.Product{}, models.ErrNotFound
	}
	row, err := s.q.ArchiveProduct(ctx, pk)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Product{}, models.ErrNotFound
	case err != nil:
		return models.Product{}, fmt.Errorf("store: archive product %s: %w", id, err)
	}
	return productFromRow(row), nil
}

// RestoreProduct brings an archived product back into the shop.
func (s *Store) RestoreProduct(ctx context.Context, id string) (models.Product, error) {
	pk, err := parseID(id)
	if err != nil {
		return models.Product{}, models.ErrNotFound
	}
	row, err := s.q.RestoreProduct(ctx, pk)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Product{}, models.ErrNotFound
	case err != nil:
		return models.Product{}, fmt.Errorf("store: restore product %s: %w", id, err)
	}
	return productFromRow(row), nil
}

// EditProduct changes a product inside a transaction that locks its row:
// edit gets the product as it is right now, and nothing can change it until
// the change is saved. A checkout that takes stock meanwhile waits, so it
// can't be undone by this write, and this write can't be based on stale stock.
func (s *Store) EditProduct(ctx context.Context, id string, edit func(*models.Product) error) (models.Product, error) {
	pk, err := parseID(id)
	if err != nil {
		return models.Product{}, models.ErrNotFound
	}
	var product models.Product
	err = s.withTx(ctx, func(_ pgx.Tx, q *sqlc.Queries) error {
		row, err := q.GetProductForUpdate(ctx, pk)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("store: edit product: %w", err)
		}
		product = productFromRow(row)
		if err := edit(&product); err != nil {
			return err
		}
		if err := product.Validate(); err != nil {
			return err
		}
		row, err = q.UpdateProduct(ctx, sqlc.UpdateProductParams{
			ID: product.ID, Name: product.Name, Description: product.Description,
			PriceCents: product.PriceCents, Stock: product.Stock, ImageUrl: product.ImageUrl,
			ArchivedAt: timePtrToTimestamptz(product.ArchivedAt),
		})
		if err != nil {
			return fmt.Errorf("store: edit product: %w", err)
		}
		product = productFromRow(row)
		return nil
	})
	return product, err
}

// ListCatalogProducts lists the products customers can buy, newest first,
// in one category if categoryID isn't nil.
func (s *Store) ListCatalogProducts(ctx context.Context, categoryID *int64, limit, offset int32) ([]models.Product, error) {
	rows, err := s.q.ListCatalogProducts(ctx, sqlc.ListCatalogProductsParams{
		CategoryID: categoryID,
		RowLimit:   limit,
		RowOffset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("store: list catalog products: %w", err)
	}
	out := make([]models.Product, 0, len(rows))
	for _, r := range rows {
		out = append(out, productFromRow(r))
	}
	return out, nil
}

// CountCatalogProducts counts the products customers can buy, in one category
// if categoryID isn't nil: the total behind the pages.
func (s *Store) CountCatalogProducts(ctx context.Context, categoryID *int64) (int64, error) {
	n, err := s.q.CountCatalogProducts(ctx, categoryID)
	if err != nil {
		return 0, fmt.Errorf("store: count catalog products: %w", err)
	}
	return n, nil
}

// GetCatalogProduct finds a product the shop sells. An archived one is not found.
func (s *Store) GetCatalogProduct(ctx context.Context, id string) (models.Product, error) {
	pk, err := parseID(id)
	if err != nil {
		return models.Product{}, models.ErrNotFound
	}
	row, err := s.q.GetCatalogProduct(ctx, pk)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Product{}, models.ErrNotFound
	case err != nil:
		return models.Product{}, fmt.Errorf("store: get catalog product %s: %w", id, err)
	}
	return productFromRow(row), nil
}
