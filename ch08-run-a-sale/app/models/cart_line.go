package models

// CartLine is one line of a cart as the customer sees it: a cart item
// joined with its product. It has no table of its own.
type CartLine struct {
	ItemID     int64
	ProductID  int64
	Name       string
	ImageUrl   string
	PriceCents int32
	Quantity   int32
}
