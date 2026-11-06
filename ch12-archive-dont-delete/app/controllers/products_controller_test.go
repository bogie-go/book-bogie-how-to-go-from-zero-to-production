package controllers_test

import (
	"fmt"
	"net/http"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("GET /products", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodGet),
			HaveField("Path", "/products"),
		)))
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
	It("lists the catalog", func() {
		factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Name = "Bogie Notebook" })

		app.Get("/products").Do().
			Status(http.StatusOK).
			Contains(`"name":"Bogie Notebook"`)
	})
})

var _ = Describe("GET /products/:id", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodGet),
			HaveField("Path", "/products/:id"),
		)))
	})

	It("answers with the product", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Name = "Bogie Ceramic Mug" })

		app.Get("/products/" + strconv.FormatInt(mug.ID, 10)).Do().
			Status(http.StatusOK).
			Contains(`"name":"Bogie Ceramic Mug"`)
	})

	It("answers 404 for a product that doesn't exist", func() {
		app.Get("/products/999999999").Do().
			Status(http.StatusNotFound).
			Contains(`"error":"no such product"`)
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

var _ = Describe("An archived product", func() {
	It("is gone from the shop", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Name = "Bogie Ceramic Mug" })
		_, err := app.Store.ArchiveProduct(GinkgoT().Context(), fmt.Sprint(mug.ID))
		Expect(err).NotTo(HaveOccurred())

		app.Get(fmt.Sprintf("/products/%d", mug.ID)).Do().Status(http.StatusNotFound)
		app.Get("/products").Do().Status(http.StatusOK).Contains(`"products":[]`)
	})
})
