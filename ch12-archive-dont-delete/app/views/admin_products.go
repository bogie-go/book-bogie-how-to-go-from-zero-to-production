package views

import (
	"time"
)

// AdminProductRequest is the body of POST /admin/products and PATCH /admin/products/:id.
type AdminProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int32  `json:"price_cents"`
	Stock       int32  `json:"stock"`
	ImageUrl    string `json:"image_url"`
}

// AdminProductResponse is one product on the wire.
type AdminProductResponse struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	PriceCents  int32      `json:"price_cents"`
	Stock       int32      `json:"stock"`
	ImageUrl    string     `json:"image_url"`
	ArchivedAt  *time.Time `json:"archived_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AdminProductsResponse is a listing.
type AdminProductsResponse struct {
	Products []AdminProductResponse `json:"products"`
}
