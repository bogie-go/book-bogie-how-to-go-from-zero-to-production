package main

import (
	"cmp"
	"fmt"
	"slices"

	"depot-playground/catalog"
)

func cheaperThan(limit int) func(catalog.Product) bool {
	return func(p catalog.Product) bool {
		return p.PriceCents < limit
	}
}

func main() {
	products := []catalog.Product{
		{Name: "Bogie Desk Mat", PriceCents: 2400, Stock: 5},
		{Name: "Bogie Sticker", PriceCents: 150, Stock: 0},
		{Name: "Bogie Mug", PriceCents: 1250, Stock: 3},
	}

	slices.SortFunc(products, func(a, b catalog.Product) int {
		return cmp.Compare(a.PriceCents, b.PriceCents)
	})

	for _, p := range products {
		fmt.Println(p.Name)
	}

	underTen := cheaperThan(1000)

	for _, p := range products {
		if underTen(p) {
			fmt.Println(p.Name)
		}
	}
}
