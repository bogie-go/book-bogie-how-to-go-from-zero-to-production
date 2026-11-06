package views

import (
	"time"
)

// AdminCouponRequest is the body of POST /admin/coupons and PATCH /admin/coupons/:id.
type AdminCouponRequest struct {
	Code        string `json:"code"`
	PriceRuleID int64  `json:"price_rule_id"`
}

// AdminCouponResponse is one coupon on the wire.
type AdminCouponResponse struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	PriceRuleID int64     `json:"price_rule_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AdminCouponsResponse is a listing.
type AdminCouponsResponse struct {
	Coupons []AdminCouponResponse `json:"coupons"`
}
