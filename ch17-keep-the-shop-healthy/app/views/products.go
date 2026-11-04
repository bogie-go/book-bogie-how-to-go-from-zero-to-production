package views

// ProductResponse is one product as the shop shows it.
type ProductResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int32  `json:"price_cents"`
	Stock       int32  `json:"stock"`
	ImageUrl    string `json:"image_url"`
}

// ProductsResponse is the catalog.
type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Page     int               `json:"page"`
	PerPage  int               `json:"per_page"`
	Total    int64             `json:"total"`
}

// CategoryNameResponse is a category as a product's page names it.
type CategoryNameResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ProductDetailResponse is one product, with the categories it's in.
type ProductDetailResponse struct {
	ProductResponse
	Categories []CategoryNameResponse `json:"categories"`
}
