package main

import (
	"fmt"

	"depot/app/models"
)

func main() {
	mug := models.Product{Name: "Bogie Ceramic Mug", PriceCents: 1400}
	fmt.Println(mug.Name, mug.PriceCents, mug.Stock)
	fmt.Println(mug.Validate())
}
