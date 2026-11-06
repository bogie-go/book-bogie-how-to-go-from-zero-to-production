package models

import (
	"fmt"
	"time"
)

// Category is a category. Imports nothing third-party: ids are int64s here and
// become bigint at the store boundary.
type Category struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate is the model's validations. Wrap every
// failure in ErrInvalid so a controller can answer 422 without knowing the
// rule:
//
//	return fmt.Errorf("%w: title is required", ErrInvalid)
func (c Category) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalid)
	}
	return nil
}

// CategorySummary is a category with the number of products customers can
// buy in it.
type CategorySummary struct {
	ID           int64
	Name         string
	ProductCount int64
}
