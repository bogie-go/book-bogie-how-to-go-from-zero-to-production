package views

import (
	"time"
)

// AdminCategoryRequest is the body of POST /admin/categories and PATCH /admin/categories/:id.
type AdminCategoryRequest struct {
	Name string `json:"name"`
}

// AdminCategoryResponse is one category on the wire.
type AdminCategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AdminCategoriesResponse is a listing.
type AdminCategoriesResponse struct {
	Categories []AdminCategoryResponse `json:"categories"`
}
