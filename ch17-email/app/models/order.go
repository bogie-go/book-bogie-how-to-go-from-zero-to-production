package models

import (
	"fmt"
	"time"
)

// OrderPlaced is the status of an order that has been placed.
const OrderPlaced = "placed"

// Order is an order. Imports nothing third-party: ids are int64s here and
// become bigint at the store boundary.
type Order struct {
	ID            int64
	UserID        int64
	Status        string
	SubtotalCents int32
	DiscountCents int32
	TotalCents    int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Validate is the model's validations. Wrap every
// failure in ErrInvalid so a controller can answer 422 without knowing the
// rule:
//
//	return fmt.Errorf("%w: title is required", ErrInvalid)
func (o Order) Validate() error {
	if o.Status != OrderPlaced {
		return fmt.Errorf("%w: status must be %q", ErrInvalid, OrderPlaced)
	}
	if o.SubtotalCents < 0 || o.DiscountCents < 0 || o.TotalCents != o.SubtotalCents-o.DiscountCents {
		return fmt.Errorf("%w: the order's totals don't add up", ErrInvalid)
	}
	return nil
}
