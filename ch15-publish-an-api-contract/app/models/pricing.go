package models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PriceContext is what flows through the pricing pipeline: the cart's lines,
// what they cost before any rule, the price so far, and the items the rules
// have given away so far.
type PriceContext struct {
	Lines         []CartLine
	OriginalCents int32
	CurrentCents  int32
	FreeItems     []FreeItem
}

// FreeItem is something a rule adds to the cart at no cost.
type FreeItem struct {
	ProductID int64
	Name      string
	ImageUrl  string
	Quantity  int32
	// Rule is the name of the rule that gave it.
	Rule string
}

// Rule is one step in the pricing pipeline. Apply takes the cart as it
// stands, in ctx, and returns it as the rule leaves it. Each kind of rule is
// its own struct, and knows which of its own settings make sense.
type Rule interface {
	Apply(ctx PriceContext) PriceContext
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

func (r PercentageDiscount) Apply(ctx PriceContext) PriceContext {
	discount := int32(int64(ctx.CurrentCents) * int64(r.Percentage) / 100)
	ctx.CurrentCents -= discount
	return ctx
}

func (r PercentageDiscount) Validate() error {
	if r.Percentage < 1 || r.Percentage > 100 {
		return fmt.Errorf("percentage must be between 1 and 100")
	}
	return nil
}

// BuyXGetY gives items away: for every Buy of one product in the cart, the
// customer gets Free more of it, at no cost. It adds them to the free items,
// and leaves the price alone.
type BuyXGetY struct {
	ProductID int64 `json:"product_id"`
	Buy       int32 `json:"buy"`
	Free      int32 `json:"free"`
}

func (r BuyXGetY) Apply(ctx PriceContext) PriceContext {
	for _, l := range ctx.Lines {
		if l.ProductID != r.ProductID {
			continue
		}
		if free := l.Quantity / r.Buy * r.Free; free > 0 {
			ctx.FreeItems = append(ctx.FreeItems, FreeItem{
				ProductID: l.ProductID, Name: l.Name, ImageUrl: l.ImageUrl, Quantity: free,
			})
		}
		break
	}
	return ctx
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

func (r MinimumPurchase) Apply(ctx PriceContext) PriceContext {
	if ctx.CurrentCents >= r.MinimumCents {
		ctx.CurrentCents -= r.DiscountCents
	}
	return ctx
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

// PricedCart is a cart's price, worked out rule by rule, and what the rules
// gave away.
type PricedCart struct {
	SubtotalCents int32
	TotalCents    int32
	Applied       []AppliedRule
	FreeItems     []FreeItem
	// Skipped holds an error for each rule that wouldn't load.
	Skipped []error
}

// DiscountCents is everything the rules took off, together.
func (p PricedCart) DiscountCents() int32 {
	return p.SubtotalCents - p.TotalCents
}

// PriceCart passes a cart through its rules, in order: each rule starts from
// the cart the one before it left. No rule can take the price below zero or
// raise it, only rules that lowered it are recorded as steps, every free item
// is named after the rule that gave it, and a rule that won't load is
// skipped, never fatal.
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
		next := rule.Apply(ctx)
		next.CurrentCents = min(max(next.CurrentCents, 0), ctx.CurrentCents)
		if next.CurrentCents < ctx.CurrentCents {
			priced.Applied = append(priced.Applied, AppliedRule{Name: pr.Name, BeforeCents: ctx.CurrentCents, AfterCents: next.CurrentCents})
		}
		for i := len(ctx.FreeItems); i < len(next.FreeItems); i++ {
			next.FreeItems[i].Rule = pr.Name
		}
		ctx = next
	}
	priced.TotalCents = ctx.CurrentCents
	priced.FreeItems = ctx.FreeItems
	return priced
}
