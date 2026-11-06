package views

// CartLineResponse is one line of the cart on the wire.
type CartLineResponse struct {
	ID             int64  `json:"id"`
	ProductID      int64  `json:"product_id"`
	Name           string `json:"name"`
	ImageUrl       string `json:"image_url"`
	UnitPriceCents int32  `json:"unit_price_cents"`
	Quantity       int32  `json:"quantity"`
	SubtotalCents  int32  `json:"subtotal_cents"`
}

// AppliedRuleResponse is one step of the cart's pricing: a rule, and the
// price before and after it.
type AppliedRuleResponse struct {
	Name             string `json:"name"`
	PriceBeforeCents int32  `json:"price_before_cents"`
	PriceAfterCents  int32  `json:"price_after_cents"`
	DiscountCents    int32  `json:"discount_cents"`
}

// CartResponse is the whole cart: every change to it answers with this.
type CartResponse struct {
	Items         []CartLineResponse    `json:"items"`
	SubtotalCents int32                 `json:"subtotal_cents"`
	DiscountCents int32                 `json:"discount_cents"`
	TotalCents    int32                 `json:"total_cents"`
	AppliedRules  []AppliedRuleResponse `json:"applied_rules"`
}

// AddCartItemRequest is the body of POST /cart_items.
type AddCartItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

// UpdateCartItemRequest is the body of PATCH /cart_items/:id.
type UpdateCartItemRequest struct {
	Quantity int32 `json:"quantity"`
}
