package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	mux.HandleFunc("GET /products/{id}", showProduct)
	mux.HandleFunc("POST /products", createProduct)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func showProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	for _, p := range products {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "product not found", http.StatusNotFound)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var p catalog.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := p.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	p.ID = len(products) + 1
	products = append(products, p)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}
