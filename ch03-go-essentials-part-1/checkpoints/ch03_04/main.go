package main

import "fmt"

func formatPrice(cents int) string {
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
}

type Product struct {
	Name       string
	PriceCents int
	Stock      int
}

func (p Product) InStock() bool {
	return p.Stock > 0
}

func main() {
	mug := Product{Name: "Bogie Mug", PriceCents: 1250, Stock: 3}
	fmt.Println(mug.Name, formatPrice(mug.PriceCents))

	if mug.InStock() {
		fmt.Println("Add to cart")
	}
}
