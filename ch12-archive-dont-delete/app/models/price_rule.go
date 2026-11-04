package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// PriceRule is a price_rule. Imports nothing third-party: ids are int64s here and
// become bigint at the store boundary.
type PriceRule struct {
	ID         int64
	Name       string
	Kind       string
	Properties json.RawMessage
	Priority   int32
	StartsAt   time.Time
	EndsAt     time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Validate is the model's validations. Wrap every
// failure in ErrInvalid so a controller can answer 422 without knowing the
// rule:
//
//	return fmt.Errorf("%w: title is required", ErrInvalid)
func (r PriceRule) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalid)
	}
	if !r.EndsAt.After(r.StartsAt) {
		return fmt.Errorf("%w: a rule must end after it starts", ErrInvalid)
	}
	if _, err := r.Rule(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	return nil
}
