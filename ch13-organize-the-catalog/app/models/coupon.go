package models

import (
	"fmt"
	"strings"
	"time"
)

// Coupon is a coupon. Imports nothing third-party: ids are int64s here and
// become bigint at the store boundary.
type Coupon struct {
	ID          int64
	Code        string
	PriceRuleID int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate is the model's validations. Wrap every
// failure in ErrInvalid so a controller can answer 422 without knowing the
// rule:
//
//	return fmt.Errorf("%w: title is required", ErrInvalid)
func (c Coupon) Validate() error {
	if c.Code == "" || c.Code != NormalizeCouponCode(c.Code) {
		return fmt.Errorf("%w: a coupon code must be upper case, with no spaces around it", ErrInvalid)
	}
	return nil
}

// NormalizeCouponCode is how codes are compared: trimmed, in upper case.
func NormalizeCouponCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
