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

func (p *Product) Restock(n int) {
	p.Stock += n
}

func main() {
	products := []Product{
		{Name: "Bogie Mug", PriceCents: 1250, Stock: 3},
		{Name: "Bogie Sticker", PriceCents: 150, Stock: 0},
	}

	products = append(products, Product{Name: "Bogie Desk Mat", PriceCents: 2400, Stock: 5})

	fmt.Println(len(products))
	fmt.Println(products[0].Name)

	for i, p := range products {
		fmt.Println(i, p.Name, p.InStock())
	}

	// p is a copy of each element, so this leaves the slice unchanged.
	for _, p := range products {
		p.Restock(10)
	}

	// Use the index to change the element itself.
	for i := range products {
		products[i].Restock(10)
	}
	fmt.Println(products[0].Stock)
}
