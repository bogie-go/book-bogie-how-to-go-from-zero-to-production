// Package order_confirmation emails a customer the receipt for their order.
// ctx first; it knows nothing about HTTP or queues, so the job and the test
// call the same function with the same values.
package order_confirmation

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"depot/app/models"
	"depot/app/services/mailer"
)

type store interface {
	GetOrder(ctx context.Context, id string) (models.Order, error)
	OrderLines(ctx context.Context, orderID int64) ([]models.OrderItem, error)
	GetUser(ctx context.Context, id string) (models.User, error)
}

type sender interface {
	Send(ctx context.Context, msg mailer.Message) error
}

// Service is order_confirmation.
type Service struct {
	store  store
	sender sender
	log    *slog.Logger
}

// New builds the service.
func New(store store, sender sender, log *slog.Logger) *Service {
	return &Service{store: store, sender: sender, log: log}
}

// Send emails the receipt for one order to the customer who placed it.
func (s *Service) Send(ctx context.Context, orderID int64) error {
	order, err := s.store.GetOrder(ctx, fmt.Sprint(orderID))
	if err != nil {
		return err
	}
	lines, err := s.store.OrderLines(ctx, order.ID)
	if err != nil {
		return err
	}
	customer, err := s.store.GetUser(ctx, fmt.Sprint(order.UserID))
	if err != nil {
		return err
	}
	return s.sender.Send(ctx, mailer.Message{
		To:      customer.Email,
		Subject: fmt.Sprintf("Your Bogie Depot order #%d", order.ID),
		Body:    receipt(order, lines),
	})
}

// receipt is the email's text: every line, then the totals.
func receipt(order models.Order, lines []models.OrderItem) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Thanks for your order! Here's your receipt.\n\n")
	for _, l := range lines {
		fmt.Fprintf(&b, "%d x %s  %s\n", l.Quantity, l.Name, money(l.UnitPriceCents*l.Quantity))
	}
	fmt.Fprintf(&b, "\nSubtotal  %s\n", money(order.SubtotalCents))
	if order.DiscountCents > 0 {
		fmt.Fprintf(&b, "Discount  -%s\n", money(order.DiscountCents))
	}
	fmt.Fprintf(&b, "Total     %s\n", money(order.TotalCents))
	return b.String()
}

// money formats cents as dollars: 2800 is $28.00.
func money(cents int32) string {
	return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
}
