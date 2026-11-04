package main

import "fmt"

func main() {
	var none []string
	names := make([]string, 0, 2)
	names = append(names, "Bogie Ceramic Mug", "Bogie Logo Sticker")
	for i, name := range names {
		fmt.Println(i, name)
	}
	fmt.Println("names has", len(names), "items; names is nil:", names == nil)
	fmt.Println("none has", len(none), "items; none is nil:", none == nil)
}
