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
}
