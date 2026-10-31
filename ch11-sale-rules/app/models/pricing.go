package models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Rule is one kind of price rule: given a cart's lines, it works out how many
// cents to take off. Each kind of rule is its own struct, and knows which of
// its own settings make sense.
type Rule interface {
	Discount(lines []CartLine) int32
	Validate() error
}

// subtotal is what the lines cost before any discount.
func subtotal(lines []CartLine) int32 {
	var total int32
	for _, l := range lines {
		total += l.PriceCents * l.Quantity
	}
	return total
}

// PercentOff takes a percentage off the whole cart.
type PercentOff struct {
	Percent int32 `json:"percent"`
}

func (r PercentOff) Discount(lines []CartLine) int32 {
	return int32(int64(subtotal(lines)) * int64(r.Percent) / 100)
}

func (r PercentOff) Validate() error {
	if r.Percent < 1 || r.Percent > 100 {
		return fmt.Errorf("percent must be between 1 and 100")
	}
	return nil
}

// BuyXGetY makes some items of one product free: for every Buy of them,
// Free more cost nothing.
type BuyXGetY struct {
	ProductID int64 `json:"product_id"`
	Buy       int32 `json:"buy"`
	Free      int32 `json:"free"`
}

func (r BuyXGetY) Validate() error {
	if r.ProductID <= 0 || r.Buy < 1 || r.Free < 1 {
		return fmt.Errorf("needs a product_id, and buy and free of at least 1")
	}
	return nil
}

func (r BuyXGetY) Discount(lines []CartLine) int32 {
	if r.Buy <= 0 || r.Free <= 0 {
		return 0
	}
	for _, l := range lines {
		if l.ProductID == r.ProductID {
			groups := l.Quantity / (r.Buy + r.Free)
			return groups * r.Free * l.PriceCents
		}
	}
	return 0
}

// AmountOff takes a fixed amount off a cart whose subtotal reaches a minimum.
type AmountOff struct {
	AmountCents      int32 `json:"amount_cents"`
	MinSubtotalCents int32 `json:"min_subtotal_cents"`
}

func (r AmountOff) Discount(lines []CartLine) int32 {
	if subtotal(lines) < r.MinSubtotalCents {
		return 0
	}
	return r.AmountCents
}

func (r AmountOff) Validate() error {
	if r.AmountCents < 1 || r.MinSubtotalCents < 0 {
		return fmt.Errorf("needs an amount_cents of at least 1, and a min_subtotal_cents of 0 or more")
	}
	return nil
}

// Rule turns a price rule's row into the rule its kind names, with the row's
// properties decoded into that rule's settings and checked.
func (r PriceRule) Rule() (Rule, error) {
	var rule Rule
	switch r.Kind {
	case "percent_off":
		rule = &PercentOff{}
	case "buy_x_get_y":
		rule = &BuyXGetY{}
	case "amount_off":
		rule = &AmountOff{}
	default:
		return nil, fmt.Errorf("unknown kind %q", r.Kind)
	}
	if err := decodeStrict(r.Properties, rule); err != nil {
		return nil, err
	}
	if err := rule.Validate(); err != nil {
		return nil, err
	}
	return rule, nil
}

// decodeStrict decodes JSON into v, refusing keys v has no field for, so a
// misspelled setting is an error instead of a silent zero.
func decodeStrict(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// AppliedRule is one rule's discount on a cart.
type AppliedRule struct {
	Name          string
	DiscountCents int32
}

// PricedCart is a cart with its discounts worked out.
type PricedCart struct {
	SubtotalCents int32
	DiscountCents int32
	TotalCents    int32
	Applied       []AppliedRule
	// Skipped holds an error for each rule that wouldn't load.
	Skipped []error
}

// PriceCart applies rules to a cart's lines, in order. Each rule's discount is
// worked out from the cart itself, not from what earlier rules left; the total
// discount never exceeds the subtotal, and only rules that take something off
// are listed. A rule that won't load is skipped, never fatal.
func PriceCart(lines []CartLine, rules []PriceRule) PricedCart {
	priced := PricedCart{SubtotalCents: subtotal(lines), Applied: []AppliedRule{}}
	for _, pr := range rules {
		rule, err := pr.Rule()
		if err != nil {
			priced.Skipped = append(priced.Skipped, fmt.Errorf("price rule %q: %w", pr.Name, err))
			continue
		}
		cents := min(rule.Discount(lines), priced.SubtotalCents-priced.DiscountCents)
		if cents > 0 {
			priced.Applied = append(priced.Applied, AppliedRule{Name: pr.Name, DiscountCents: cents})
			priced.DiscountCents += cents
		}
	}
	priced.TotalCents = priced.SubtotalCents - priced.DiscountCents
	return priced
}
