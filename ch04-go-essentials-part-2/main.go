package main

import (
	"context"
	"fmt"
	"time"

	"depot-playground/catalog"
)

func restockAll(ctx context.Context, products []catalog.Product) error {
	for i := range products {
		if err := ctx.Err(); err != nil {
			return err
		}
		time.Sleep(30 * time.Millisecond) // pretend this is slow
		products[i].Restock(10)
	}
	return nil
}

func main() {
	products := []catalog.Product{
		{Name: "Bogie Mug", Stock: 3},
		{Name: "Bogie Sticker", Stock: 0},
		{Name: "Bogie Desk Mat", Stock: 5},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := restockAll(ctx, products)
	fmt.Println(err)

	for _, p := range products {
		fmt.Println(p.Name, p.Stock)
	}
}
