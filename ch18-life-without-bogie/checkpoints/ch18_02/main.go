package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
	Stock      int    `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Bogie Logo Sticker", PriceCents: 300, Stock: 500},
	{ID: 2, Name: "Bogie Ceramic Mug", PriceCents: 1400, Stock: 150},
	{ID: 3, Name: "Bogie Desk Mat", PriceCents: 2500, Stock: 100},
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
