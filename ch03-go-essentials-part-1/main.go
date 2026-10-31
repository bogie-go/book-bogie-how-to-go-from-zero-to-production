package main

import (
	"errors"
	"fmt"
)

type Product struct {
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
	Stock      int    `json:"stock"`
}

func (p Product) InStock() bool {
	return p.Stock > 0
}

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

var ErrOutOfStock = errors.New("out of stock")

func (p *Product) Sell(qty int) error {
	if qty > p.Stock {
		return fmt.Errorf("selling %d of %q: %w", qty, p.Name, ErrOutOfStock)
	}
	p.Stock -= qty
	return nil
}

func main() {
	mug := Product{Name: "Bogie Mug", PriceCents: 1250, Stock: 3}

	err := mug.Sell(10)
	if errors.Is(err, ErrOutOfStock) {
		fmt.Println("Sorry:", err)
	}
}
