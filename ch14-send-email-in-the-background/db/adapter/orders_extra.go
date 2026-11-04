package adapter

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"

	"depot/app/models"
	"depot/db/sqlc"
)

// AfterOrder runs inside PlaceOrder's transaction, once the order is written:
// whatever it does commits or rolls back with the order. It may be nil.
type AfterOrder func(ctx context.Context, tx pgx.Tx, order models.Order) error

// PlaceOrder turns a customer's cart into an order, in one transaction: it
// takes every line out of stock, writes the order and its lines with today's
// prices, and empties the cart. If any step fails, none of them happened.
func (s *Store) PlaceOrder(ctx context.Context, userID, cartID int64, rules []models.PriceRule, after AfterOrder) (models.Order, error) {
	var order models.Order
	err := s.withTx(ctx, func(tx pgx.Tx, q *sqlc.Queries) error {
		if _, err := q.LockCart(ctx, cartID); err != nil {
			return fmt.Errorf("store: lock cart: %w", err)
		}
		lines, err := cartLines(ctx, q, cartID)
		if err != nil {
			return err
		}
		if len(lines) == 0 {
			return fmt.Errorf("%w: the cart is empty", models.ErrInvalid)
		}

		byProduct := slices.Clone(lines)
		slices.SortFunc(byProduct, func(a, b models.CartLine) int { return cmp.Compare(a.ProductID, b.ProductID) })
		for _, l := range byProduct {
			n, err := q.TakeStock(ctx, sqlc.TakeStockParams{ID: l.ProductID, Quantity: l.Quantity})
			if err != nil {
				return fmt.Errorf("store: take stock: %w", err)
			}
			if n == 0 {
				return fmt.Errorf("%w: %s is no longer available", models.ErrInvalid, l.Name)
			}
		}

		priced := models.PriceCart(lines, rules)
		order = models.Order{
			UserID: userID, Status: models.OrderPlaced,
			SubtotalCents: priced.SubtotalCents, DiscountCents: priced.DiscountCents, TotalCents: priced.TotalCents,
		}
		if err := order.Validate(); err != nil {
			return err
		}
		row, err := q.CreateOrder(ctx, sqlc.CreateOrderParams{
			UserID: order.UserID, Status: order.Status,
			SubtotalCents: order.SubtotalCents, DiscountCents: order.DiscountCents, TotalCents: order.TotalCents,
		})
		if err != nil {
			return fmt.Errorf("store: create order: %w", err)
		}
		order = orderFromRow(row)

		for _, l := range lines {
			item := models.OrderItem{OrderID: order.ID, ProductID: l.ProductID, Name: l.Name, Quantity: l.Quantity, UnitPriceCents: l.PriceCents}
			if err := item.Validate(); err != nil {
				return err
			}
			if _, err := q.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
				OrderID: item.OrderID, ProductID: item.ProductID, Name: item.Name,
				Quantity: item.Quantity, UnitPriceCents: item.UnitPriceCents,
			}); err != nil {
				return fmt.Errorf("store: create order item: %w", err)
			}
		}

		if _, err := q.DeleteCartItemsForCart(ctx, cartID); err != nil {
			return fmt.Errorf("store: empty cart: %w", err)
		}
		if after != nil {
			if err := after(ctx, tx, order); err != nil {
				return fmt.Errorf("store: after order: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

// OrderLines lists an order's lines, in the order they were written.
func (s *Store) OrderLines(ctx context.Context, orderID int64) ([]models.OrderItem, error) {
	rows, err := s.q.ListOrderItemsForOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("store: order lines: %w", err)
	}
	items := make([]models.OrderItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, orderItemFromRow(r))
	}
	return items, nil
}

// FindOrderForUser finds an order within one customer's own orders. Another
// customer's order is simply not found.
func (s *Store) FindOrderForUser(ctx context.Context, id string, userID int64) (models.Order, error) {
	pk, err := parseID(id)
	if err != nil {
		return models.Order{}, models.ErrNotFound
	}
	row, err := s.q.GetOrderForUser(ctx, sqlc.GetOrderForUserParams{ID: pk, UserID: userID})
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Order{}, models.ErrNotFound
	}
	if err != nil {
		return models.Order{}, fmt.Errorf("store: find order: %w", err)
	}
	return orderFromRow(row), nil
}
