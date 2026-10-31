package main

import (
	"encoding/json"
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

func main() {
	mug := Product{Name: "Bogie Mug", PriceCents: 1250, Stock: 3}

	data, err := json.Marshal(mug)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}
