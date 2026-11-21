package models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PriceContext is what flows through the pricing pipeline: the cart's lines,
// what they cost before any rule, the price so far, and what the rules have
// done so far: the steps they recorded, and the items they gave away.
type PriceContext struct {
	Lines         []CartLine
	OriginalCents int32
	CurrentCents  int32
	AppliedRules  []AppliedRule
	FreeItems     []FreeItem
}

// AppliedRule is what one rule did: its name, and the price before and after
// it.
type AppliedRule struct {
	Name        string
	BeforeCents int32
	AfterCents  int32
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

// RuleInfo is what every kind of rule carries from its row.
type RuleInfo struct {
	Name     string
	Priority int32
}

// Info returns the rule's name and priority.
func (i RuleInfo) Info() RuleInfo { return i }

// Rule is one step in the pricing pipeline. Apply takes the cart as it
// stands, in ctx, decides whether the rule applies, and returns the cart as
// the rule leaves it, with a note of what it did. Each kind of rule is its own
// struct, and knows which of its own settings make sense.
type Rule interface {
	Apply(ctx PriceContext) PriceContext
	Validate() error
	Info() RuleInfo
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
	RuleInfo   `json:"-"`
	Percentage int32 `json:"percentage"`
}

func (r PercentageDiscount) Apply(ctx PriceContext) PriceContext {
	discount := int32(int64(ctx.CurrentCents) * int64(r.Percentage) / 100)
	if discount == 0 {
		return ctx
	}
	before := ctx.CurrentCents
	ctx.CurrentCents -= discount
	ctx.AppliedRules = append(ctx.AppliedRules,
		AppliedRule{Name: r.Name, BeforeCents: before, AfterCents: ctx.CurrentCents})
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
	RuleInfo  `json:"-"`
	ProductID int64 `json:"product_id"`
	Buy       int32 `json:"buy"`
	Free      int32 `json:"free"`
}

func (r BuyXGetY) Apply(ctx PriceContext) PriceContext {
	for _, l := range ctx.Lines {
		if l.ProductID != r.ProductID {
			continue
		}
		free := l.Quantity / r.Buy * r.Free
		if free == 0 {
			return ctx
		}
		ctx.FreeItems = append(ctx.FreeItems, FreeItem{
			ProductID: l.ProductID, Name: l.Name, ImageUrl: l.ImageUrl,
			Quantity: free, Rule: r.Name,
		})
		ctx.AppliedRules = append(ctx.AppliedRules,
			AppliedRule{Name: r.Name, BeforeCents: ctx.CurrentCents, AfterCents: ctx.CurrentCents})
		return ctx
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
	RuleInfo      `json:"-"`
	MinimumCents  int32 `json:"minimum_cents"`
	DiscountCents int32 `json:"discount_cents"`
}

func (r MinimumPurchase) Apply(ctx PriceContext) PriceContext {
	discount := min(r.DiscountCents, ctx.CurrentCents)
	if ctx.CurrentCents < r.MinimumCents || discount == 0 {
		return ctx
	}
	before := ctx.CurrentCents
	ctx.CurrentCents -= discount
	ctx.AppliedRules = append(ctx.AppliedRules,
		AppliedRule{Name: r.Name, BeforeCents: before, AfterCents: ctx.CurrentCents})
	return ctx
}

func (r MinimumPurchase) Validate() error {
	if r.DiscountCents < 1 || r.MinimumCents < 0 {
		return fmt.Errorf("needs a discount_cents of at least 1, and a minimum_cents of 0 or more")
	}
	return nil
}

// NewRule builds the rule a row's kind names, with the row's name and
// priority, and its properties decoded into that rule's settings and checked.
// It's what ActiveRecord does for you when it reads a type column.
func NewRule(pr PriceRule) (Rule, error) {
	info := RuleInfo{Name: pr.Name, Priority: pr.Priority}
	var rule Rule
	switch pr.Kind {
	case "percentage_discount":
		rule = &PercentageDiscount{RuleInfo: info}
	case "buy_x_get_y":
		rule = &BuyXGetY{RuleInfo: info}
	case "minimum_purchase":
		rule = &MinimumPurchase{RuleInfo: info}
	default:
		return nil, fmt.Errorf("unknown kind %q", pr.Kind)
	}
	if err := decodeStrict(pr.Properties, rule); err != nil {
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

// PricedCart is a cart's price, worked out rule by rule, and what the rules
// gave away.
type PricedCart struct {
	SubtotalCents int32
	TotalCents    int32
	Applied       []AppliedRule
	FreeItems     []FreeItem
}

// DiscountCents is everything the rules took off, together.
func (p PricedCart) DiscountCents() int32 {
	return p.SubtotalCents - p.TotalCents
}

// PriceCart passes a cart through its rules, in order: each rule starts from
// the cart the one before it left.
func PriceCart(lines []CartLine, rules []Rule) PricedCart {
	ctx := PriceContext{Lines: lines, OriginalCents: subtotal(lines)}
	ctx.CurrentCents = ctx.OriginalCents
	for _, rule := range rules {
		ctx = rule.Apply(ctx)
	}
	return PricedCart{
		SubtotalCents: ctx.OriginalCents,
		TotalCents:    ctx.CurrentCents,
		Applied:       ctx.AppliedRules,
		FreeItems:     ctx.FreeItems,
	}
}
