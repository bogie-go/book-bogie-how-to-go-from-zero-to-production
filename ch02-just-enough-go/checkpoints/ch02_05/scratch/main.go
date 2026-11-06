package main

import "fmt"

type Product struct {
	Name  string
	Stock int32
}

// Take removes n items from stock.
func (p Product) Take(n int32) {
	p.Stock -= n
}

func main() {
	mug := Product{Name: "Bogie Ceramic Mug", Stock: 150}
	mug.Take(3)
	fmt.Println(mug.Stock)
}
