package models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PriceContext is what flows through the pricing pipeline: the cart's lines,
// what they cost before any rule, and the price so far.
type PriceContext struct {
	Lines         []CartLine
	OriginalCents int32
	CurrentCents  int32
}

// Rule is one step in the pricing pipeline. Apply takes the price so far, in
// ctx, and returns the new price. Each kind of rule is its own struct, and
// knows which of its own settings make sense.
type Rule interface {
	Apply(ctx PriceContext) int32
	Validate() error
}

// subtotal is what the lines cost before any rule.
func subtotal(lines []CartLine) int32 {
	var total int32
	for _, l := range lines {
		total += l.PriceCents * l.Quantity
	}
	return total
}

// PercentageDiscount takes a percentage off the price so far.
type PercentageDiscount struct {
	Percentage int32 `json:"percentage"`
}

func (r PercentageDiscount) Apply(ctx PriceContext) int32 {
	discount := int32(int64(ctx.CurrentCents) * int64(r.Percentage) / 100)
	return ctx.CurrentCents - discount
}

func (r PercentageDiscount) Validate() error {
	if r.Percentage < 1 || r.Percentage > 100 {
		return fmt.Errorf("percentage must be between 1 and 100")
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

func (r BuyXGetY) Apply(ctx PriceContext) int32 {
	for _, l := range ctx.Lines {
		if l.ProductID == r.ProductID {
			groups := l.Quantity / (r.Buy + r.Free)
			return ctx.CurrentCents - groups*r.Free*l.PriceCents
		}
	}
	return ctx.CurrentCents
}

func (r BuyXGetY) Validate() error {
	if r.ProductID <= 0 || r.Buy < 1 || r.Free < 1 {
		return fmt.Errorf("needs a product_id, and buy and free of at least 1")
	}
	return nil
}

// MinimumPurchase takes a fixed amount off, once the price so far reaches
// a minimum.
type MinimumPurchase struct {
	MinimumCents  int32 `json:"minimum_cents"`
	DiscountCents int32 `json:"discount_cents"`
}

func (r MinimumPurchase) Apply(ctx PriceContext) int32 {
	if ctx.CurrentCents < r.MinimumCents {
		return ctx.CurrentCents
	}
	return ctx.CurrentCents - r.DiscountCents
}

func (r MinimumPurchase) Validate() error {
	if r.DiscountCents < 1 || r.MinimumCents < 0 {
		return fmt.Errorf("needs a discount_cents of at least 1, and a minimum_cents of 0 or more")
	}
	return nil
}

// Rule turns a price rule's row into the rule its kind names, with the row's
// properties decoded into that rule's settings and checked.
func (r PriceRule) Rule() (Rule, error) {
	var rule Rule
	switch r.Kind {
	case "percentage_discount":
		rule = &PercentageDiscount{}
	case "buy_x_get_y":
		rule = &BuyXGetY{}
	case "minimum_purchase":
		rule = &MinimumPurchase{}
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

// AppliedRule is one step of the pipeline: a rule, and the price before and
// after it.
type AppliedRule struct {
	Name        string
	BeforeCents int32
	AfterCents  int32
}

// PricedCart is a cart's price, worked out rule by rule.
type PricedCart struct {
	SubtotalCents int32
	TotalCents    int32
	Applied       []AppliedRule
	// Skipped holds an error for each rule that wouldn't load.
	Skipped []error
}

// DiscountCents is everything the rules took off, together.
func (p PricedCart) DiscountCents() int32 {
	return p.SubtotalCents - p.TotalCents
}

// PriceCart passes a cart's price through its rules, in order: each rule
// starts from the price the one before it left. No rule can take the price
// below zero or raise it, only rules that lowered it are recorded, and a rule
// that won't load is skipped, never fatal.
func PriceCart(lines []CartLine, rules []PriceRule) PricedCart {
	ctx := PriceContext{Lines: lines, OriginalCents: subtotal(lines)}
	ctx.CurrentCents = ctx.OriginalCents
	priced := PricedCart{SubtotalCents: ctx.OriginalCents, Applied: []AppliedRule{}}
	for _, pr := range rules {
		rule, err := pr.Rule()
		if err != nil {
			priced.Skipped = append(priced.Skipped, fmt.Errorf("price rule %q: %w", pr.Name, err))
			continue
		}
		after := max(rule.Apply(ctx), 0)
		if after < ctx.CurrentCents {
			priced.Applied = append(priced.Applied, AppliedRule{Name: pr.Name, BeforeCents: ctx.CurrentCents, AfterCents: after})
			ctx.CurrentCents = after
		}
	}
	priced.TotalCents = ctx.CurrentCents
	return priced
}
