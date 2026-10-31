package main

import "fmt"

type Product struct {
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
	Stock      int    `json:"stock"`
}

func (p Product) InStock() bool {
	return p.Stock > 0
}

// Restock has a pointer receiver, so it changes the original product.
func (p *Product) Restock(n int) {
	p.Stock += n
}

func main() {
	mug := Product{Name: "Bogie Mug", PriceCents: 1250, Stock: 3}
	mug.Restock(10)
	fmt.Println(mug.Stock)

	p := &mug    // & takes the address of mug; p is a *Product
	p.Stock = 20 // changes mug itself
	fmt.Println(mug.Stock)
}
