package main

import "fmt"

func main() {
	stock := map[string]int{
		"Bogie Mug":     3,
		"Bogie Sticker": 0,
	}

	stock["Bogie Desk Mat"] = 5
	fmt.Println(stock["Bogie Mug"])

	count, ok := stock["Bogie Keycap"]
	fmt.Println(count, ok)
}
