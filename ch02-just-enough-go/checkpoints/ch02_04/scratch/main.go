package main

import (
	"fmt"
	"math"
)

// formatPrice shows cents as dollars and cents.
func formatPrice(cents int32) string {
	return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
}

func main() {
	a, b := 0.1, 0.2
	fmt.Println(a + b)

	fmt.Println(formatPrice(1400), formatPrice(305))
	fmt.Println(formatPrice(math.MaxInt32))
}
