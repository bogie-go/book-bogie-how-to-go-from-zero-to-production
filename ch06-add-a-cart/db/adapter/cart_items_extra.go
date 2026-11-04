package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"depot/app/models"
	"depot/db/sqlc"
)

// CartContents lists a cart's lines with their products, oldest first.
func (s *Store) CartContents(ctx context.Context, cartID int64) ([]models.CartLine, error) {
	rows, err := s.q.ListCartItemsForCart(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("store: cart contents: %w", err)
	}
	lines := make([]models.CartLine, 0, len(rows))
	for _, r := range rows {
		lines = append(lines, models.CartLine{
			ItemID: r.ID, ProductID: r.ProductID, Name: r.Name,
			ImageUrl: r.ImageUrl, PriceCents: r.PriceCents, Quantity: r.Quantity,
		})
	}
	return lines, nil
}

// AddCartItem adds quantity of a product to a cart: to the product's existing
// line if there is one, or as a new line, in one statement. The line may not
// ask for more than the product has in stock.
func (s *Store) AddCartItem(ctx context.Context, cartID, productID int64, quantity int32) error {
	if quantity <= 0 {
		return fmt.Errorf("%w: quantity must be at least 1", models.ErrInvalid)
	}
	product, err := s.q.GetProduct(ctx, productID)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("store: add to cart: %w", err)
	}
	if quantity > product.Stock {
		return outOfStock(product.Stock)
	}
	_, err = s.q.AddToCart(ctx, sqlc.AddToCartParams{
		CartID: cartID, ProductID: productID, Quantity: quantity, Stock: product.Stock,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return outOfStock(product.Stock)
	}
	if err != nil {
		return fmt.Errorf("store: add to cart: %w", err)
	}
	return nil
}

func outOfStock(stock int32) error {
	return fmt.Errorf("%w: only %d left in stock", models.ErrInvalid, stock)
}

// SetCartItemQuantity changes a line's quantity, within the product's stock.
func (s *Store) SetCartItemQuantity(ctx context.Context, id string, quantity int32) (models.CartItem, error) {
	item, err := s.GetCartItem(ctx, id)
	if err != nil {
		return models.CartItem{}, err
	}
	product, err := s.q.GetProduct(ctx, item.ProductID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("store: set cart quantity: %w", err)
	}
	if quantity > product.Stock {
		return models.CartItem{}, outOfStock(product.Stock)
	}
	item.Quantity = quantity
	return s.UpdateCartItem(ctx, item)
}
