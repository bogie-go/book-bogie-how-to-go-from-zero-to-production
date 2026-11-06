package views

// CategoryResponse is one category on the wire.
type CategoryResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ProductCount int64  `json:"product_count"`
}

// CategoriesResponse is a listing.
type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
}
