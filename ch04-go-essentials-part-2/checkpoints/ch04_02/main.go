package main

import (
	"fmt"

	"depot-playground/catalog"
)

func main() {
	mug := catalog.Product{Name: "Bogie Mug", PriceCents: 1250, Stock: 3}
	fmt.Println(mug.Name, mug.InStock())
	fmt.Println(mug)
}
