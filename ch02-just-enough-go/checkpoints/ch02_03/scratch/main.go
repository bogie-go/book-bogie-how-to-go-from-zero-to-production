package main

import (
	"fmt"
	"strconv"
)

// parsePrice turns text like "1400", a price in cents, into a number.
func parsePrice(s string) (int32, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("price %q isn't a number", s)
	}
	return int32(n), nil
}

func main() {
	price, err := parsePrice("1400")
	fmt.Println(price, err)

	price, err = parsePrice("cheap")
	fmt.Println(price, err)
}
