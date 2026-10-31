package main

import (
	"fmt"

	"depot-playground/catalog"
)

type Priced interface {
	Total() int
}

func receiptTotal(lines []Priced) int {
	total := 0
	for _, line := range lines {
		total += line.Total()
	}
	return total
}

type CartLine struct {
	Product  catalog.Product
	Quantity int
}

func (l CartLine) Total() int {
	return l.Product.PriceCents * l.Quantity
}

type GiftWrap struct{}

func (GiftWrap) Total() int {
	return 300
}

func main() {
	mug := catalog.Product{Name: "Bogie Mug", PriceCents: 1250, Stock: 3}

	lines := []Priced{
		CartLine{Product: mug, Quantity: 2},
		GiftWrap{},
	}
	fmt.Println(receiptTotal(lines))
}
