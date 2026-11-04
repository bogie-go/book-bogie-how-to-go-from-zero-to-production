package controllers_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("POST /cart_items", func() {
	It("answers with the whole cart, adding to the line for the same product", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) {
			p.PriceCents = 1400
			p.Stock = 10
		})
		body := map[string]any{"product_id": mug.ID, "quantity": 1}

		app.Post("/cart_items", body).Do().Status(http.StatusOK)
		app.Post("/cart_items", body).Do().
			Status(http.StatusOK).
			Contains(`"quantity":2`).
			Contains(`"total_cents":2800`)
	})
})
