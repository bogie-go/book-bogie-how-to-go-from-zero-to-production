package main

import (
	"encoding/json"
	"log"
	"net/http"

	"depot-playground/catalog"
)

var products = []catalog.Product{
	{ID: 1, Name: "Bogie Mug", PriceCents: 1250, Stock: 3},
	{ID: 2, Name: "Bogie Sticker", PriceCents: 150, Stock: 0},
	{ID: 3, Name: "Bogie Desk Mat", PriceCents: 2400, Stock: 5},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", listProducts)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
