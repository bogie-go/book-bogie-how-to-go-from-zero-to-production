package models

import (
	"fmt"
	"time"
)

// Product is a product. Imports nothing third-party: ids are int64s here and
// become bigint at the store boundary.
type Product struct {
	ID          int64
	Name        string
	Description string
	PriceCents  int32
	Stock       int32
	ImageUrl    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate is the model's validations. Wrap every
// failure in ErrInvalid so a controller can answer 422 without knowing the
// rule:
//
//	return fmt.Errorf("%w: title is required", ErrInvalid)
func (p Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalid)
	}
	if p.PriceCents <= 0 {
		return fmt.Errorf("%w: price must be more than zero", ErrInvalid)
	}
	if p.Stock < 0 {
		return fmt.Errorf("%w: stock can't be negative", ErrInvalid)
	}
	return nil
}
