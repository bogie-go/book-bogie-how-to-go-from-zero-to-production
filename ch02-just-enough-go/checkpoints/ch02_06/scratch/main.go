package main

import (
	"errors"
	"fmt"
)

type Product struct {
	Name       string
	PriceCents int32
	Stock      int32
}

var ErrSoldOut = errors.New("sold out")

// Take removes n items from stock, or fails without changing anything.
func (p *Product) Take(n int32) error {
	if n > p.Stock {
		return ErrSoldOut
	}
	p.Stock -= n
	return nil
}

// Find returns the product with the name, and whether there was one.
func Find(products []Product, name string) (Product, bool) {
	for _, p := range products {
		if p.Name == name {
			return p, true
		}
	}
	return Product{}, false
}

func main() {
	products := []Product{
		{Name: "Bogie Ceramic Mug", PriceCents: 1400, Stock: 2},
		{Name: "Bogie Logo Sticker", PriceCents: 300, Stock: 500},
	}

	for i := range products {
		if err := products[i].Take(3); err != nil {
			fmt.Println(products[i].Name+":", err)
			continue
		}
		fmt.Println(products[i].Name+": took 3, left", products[i].Stock)
	}

	if p, ok := Find(products, "Bogie Logo Sticker"); ok {
		fmt.Println("found", p.Name, "at", p.PriceCents, "cents")
	}
	if _, ok := Find(products, "Bogie Keycap"); !ok {
		fmt.Println("no keycaps")
	}

	stockByName := map[string]int32{}
	for _, p := range products {
		stockByName[p.Name] = p.Stock
	}
	fmt.Println(stockByName["Bogie Ceramic Mug"], stockByName["Bogie Logo Sticker"])
}
