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

func (p Product) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.PriceCents < 0 {
		return errors.New("price can't be negative")
	}
	return nil
}

func main() {
	unnamed := Product{Name: "", PriceCents: 500}
	if err := unnamed.Validate(); err != nil {
		fmt.Println("Can't save:", err)
		return
	}
}
