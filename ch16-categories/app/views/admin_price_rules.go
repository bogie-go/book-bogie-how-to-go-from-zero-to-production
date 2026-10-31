package views

import (
	"encoding/json"
	"time"
)

// AdminPriceRuleRequest is the body of POST /admin/price_rules and PATCH /admin/price_rules/:id.
type AdminPriceRuleRequest struct {
	Name       string          `json:"name"`
	Kind       string          `json:"kind"`
	Properties json.RawMessage `json:"properties"`
	Priority   int32           `json:"priority"`
	StartsAt   time.Time       `json:"starts_at"`
	EndsAt     time.Time       `json:"ends_at"`
}

// AdminPriceRuleResponse is one price_rule on the wire.
type AdminPriceRuleResponse struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	Kind       string          `json:"kind"`
	Properties json.RawMessage `json:"properties"`
	Priority   int32           `json:"priority"`
	StartsAt   time.Time       `json:"starts_at"`
	EndsAt     time.Time       `json:"ends_at"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// AdminPriceRulesResponse is a listing.
type AdminPriceRulesResponse struct {
	PriceRules []AdminPriceRuleResponse `json:"price_rules"`
}
