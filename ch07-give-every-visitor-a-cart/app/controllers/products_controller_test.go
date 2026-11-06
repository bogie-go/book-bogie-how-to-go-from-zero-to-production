package controllers_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/app/views"
	"depot/test/factories"
)

var _ = Describe("GET /products", func() {
	Context("with no products", func() {
		It("is an empty list", func() {
			app.Get("/products").Do().Status(http.StatusOK).Contains("\"products\":[]")
		})
	})

	Context("with products", func() {
		var first, second models.Product

		BeforeEach(func() {
			first = factories.Product(GinkgoT(), app.Store)
			second = factories.Product(GinkgoT(), app.Store)
		})

		It("lists them, newest first", func() {
			var list views.ProductsResponse
			app.Get("/products").Do().Status(http.StatusOK).JSON(&list)

			Expect(list.Products).To(HaveLen(2))
			Expect(list.Products[0].ID).To(Equal(second.ID))
			Expect(list.Products[1].ID).To(Equal(first.ID))
		})
	})

	Context("with a product", func() {
		BeforeEach(func() {
			factories.Product(GinkgoT(), app.Store, func(p *models.Product) {
				p.Name = "Bogie Ceramic Mug"
				p.PriceCents = 1400
			})
		})

		It("lists it, with the field names the shop's pages read", func() {
			app.Get("/products").Do().
				Status(http.StatusOK).
				Contains(`"name":"Bogie Ceramic Mug"`).
				Contains(`"price_cents":1400`)
		})
	})
})

var _ = Describe("GET /products/:id", func() {
	Context("when the product exists", func() {
		var record models.Product

		BeforeEach(func() {
			record = factories.Product(GinkgoT(), app.Store)
		})

		It("shows it", func() {
			var shown views.ProductResponse
			app.Get("/products/" + fmt.Sprint(record.ID)).Do().Status(http.StatusOK).JSON(&shown)

			Expect(shown.ID).To(Equal(record.ID))
		})
	})

	Context("when the product does not exist", func() {
		It("is not found", func() {
			app.Get("/products/999999999").Do().Status(http.StatusNotFound)
		})
	})
})

var _ = Describe("GET /products when the database fails", func() {
	BeforeEach(func() {
		app.Store.Close()
	})

	It("answers 500 without the database's error", func() {
		res := app.Get("/products").Do().
			Status(http.StatusInternalServerError).
			Contains(`"error":"internal error"`)

		Expect(string(res.Body)).NotTo(ContainSubstring("pool"))
	})
})
