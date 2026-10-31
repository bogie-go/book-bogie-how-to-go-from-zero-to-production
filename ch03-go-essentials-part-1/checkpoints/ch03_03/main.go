package main

import "fmt"

func formatPrice(cents int) string {
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
}

func splitPrice(cents int) (int, int) {
	return cents / 100, cents % 100
}

func main() {
	fmt.Println(formatPrice(1250))

	whole, fraction := splitPrice(1250)
	fmt.Println(whole, fraction) // 12 50

	onlyWhole, _ := splitPrice(1250)
	fmt.Println(onlyWhole)
}
