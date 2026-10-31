package catalog

import (
	"errors"
	"fmt"
)

var ErrOutOfStock = errors.New("out of stock")

type Product struct {
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
	Stock      int    `json:"stock"`
}

func (p Product) InStock() bool {
	return p.Stock > 0
}

// Restock, Validate, and Sell move here unchanged from Chapter 3.

func (p *Product) Restock(n int) {
	p.Stock += n
}

func (p Product) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.PriceCents < 0 {
		return errors.New("price can't be negative")
	}
	return nil
}

func (p *Product) Sell(qty int) error {
	if qty > p.Stock {
		return fmt.Errorf("selling %d of %q: %w", qty, p.Name, ErrOutOfStock)
	}
	p.Stock -= qty
	return nil
}
