package models

import (
	"fmt"
	"time"
)

// OrderItem is an order_item. Imports nothing third-party: ids are int64s here and
// become bigint at the store boundary.
type OrderItem struct {
	ID             int64
	OrderID        int64
	ProductID      int64
	Name           string
	Quantity       int32
	UnitPriceCents int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Validate is the model's validations. Wrap every
// failure in ErrInvalid so a controller can answer 422 without knowing the
// rule:
//
//	return fmt.Errorf("%w: title is required", ErrInvalid)
func (i OrderItem) Validate() error {
	if i.Name == "" || i.Quantity < 1 || i.UnitPriceCents < 0 {
		return fmt.Errorf("%w: an order line needs a name, a quantity of at least 1, and a price", ErrInvalid)
	}
	return nil
}
