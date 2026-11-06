package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
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

func (p Product) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.PriceCents <= 0 {
		return errors.New("price must be more than zero")
	}
	return nil
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
	writeJSON(w, http.StatusOK, products)
}

func showProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	for _, p := range products {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}

	writeError(w, http.StatusNotFound, "product not found")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := p.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	p.ID = len(products) + 1
	products = append(products, p)
	writeJSON(w, http.StatusCreated, p)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
