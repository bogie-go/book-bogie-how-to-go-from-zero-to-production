package adapter

import (
	"context"
	"fmt"

	"depot/app/models"
)

// CategoriesForProduct lists a product's categories by name: the Bogie version
// of `product.categories`. One query, run only when the code calls this.
func (s *Store) CategoriesForProduct(ctx context.Context, productID int64) ([]models.Category, error) {
	rows, err := s.q.ListCategoriesForProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("store: categories for product: %w", err)
	}
	out := make([]models.Category, 0, len(rows))
	for _, row := range rows {
		out = append(out, categoryFromRow(row))
	}
	return out, nil
}

// ListCategorySummaries lists every category by name, with how many products
// customers can buy in each.
func (s *Store) ListCategorySummaries(ctx context.Context) ([]models.CategorySummary, error) {
	rows, err := s.q.ListCategorySummaries(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: list category summaries: %w", err)
	}
	out := make([]models.CategorySummary, 0, len(rows))
	for _, r := range rows {
		out = append(out, models.CategorySummary{ID: r.ID, Name: r.Name, ProductCount: r.ProductCount})
	}
	return out, nil
}
