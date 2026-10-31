package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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
	log.Fatal(http.ListenAndServe(":8080", logRequests(mux)))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
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
	var p catalog.Product
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

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
