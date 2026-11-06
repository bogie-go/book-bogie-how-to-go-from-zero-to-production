package views

import "time"

// PlaceOrderRequest is the body of POST /orders.
type PlaceOrderRequest struct {
	CouponCode string `json:"coupon_code"`
}

// OrderItemResponse is one line of an order, at the price it was bought for.
type OrderItemResponse struct {
	ProductID      int64  `json:"product_id"`
	Name           string `json:"name"`
	Quantity       int32  `json:"quantity"`
	UnitPriceCents int32  `json:"unit_price_cents"`
	SubtotalCents  int32  `json:"subtotal_cents"`
}

// OrderResponse is an order, as its receipt shows it.
type OrderResponse struct {
	ID            int64               `json:"id"`
	Status        string              `json:"status"`
	SubtotalCents int32               `json:"subtotal_cents"`
	DiscountCents int32               `json:"discount_cents"`
	TotalCents    int32               `json:"total_cents"`
	Items         []OrderItemResponse `json:"items"`
	CreatedAt     time.Time           `json:"created_at"`
}
