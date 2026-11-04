package controllers_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/app/views"
	"depot/test/factories"
)

var _ = Describe("GET /products", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodGet),
			HaveField("Path", "/products"),
		)))
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

var _ = Describe("Products in categories", func() {
	var stickers, desk models.Category
	var sticker models.Product

	BeforeEach(func() {
		stickers = factories.Category(GinkgoT(), app.Store, func(c *models.Category) { c.Name = "Stickers" })
		desk = factories.Category(GinkgoT(), app.Store, func(c *models.Category) { c.Name = "Desk" })
		sticker = factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Name = "Logo Sticker" })
		mat := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Name = "Desk Mat" })
		Expect(app.Store.SetProductCategories(GinkgoT().Context(), fmt.Sprint(sticker.ID), []int64{stickers.ID, desk.ID})).To(Succeed())
		Expect(app.Store.SetProductCategories(GinkgoT().Context(), fmt.Sprint(mat.ID), []int64{desk.ID})).To(Succeed())
	})

	It("lists one category's products, and only those", func() {
		res := app.Get(fmt.Sprintf("/products?category_id=%d", stickers.ID)).Do().Status(http.StatusOK)
		Expect(string(res.Body)).To(ContainSubstring("Logo Sticker"))
		Expect(string(res.Body)).NotTo(ContainSubstring("Desk Mat"))
	})

	It("lists a product once, however many categories it's in", func() {
		res := app.Get(fmt.Sprintf("/products?category_id=%d", desk.ID)).Do().Status(http.StatusOK)
		Expect(strings.Count(string(res.Body), "Logo Sticker")).To(Equal(1))
	})

	It("shows a product's categories on its own page", func() {
		app.Get(fmt.Sprintf("/products/%d", sticker.ID)).Do().
			Status(http.StatusOK).
			Contains(`"categories":[{"id":`).
			Contains(`"name":"Desk"`).
			Contains(`"name":"Stickers"`)
	})

	It("refuses a category_id that isn't a number", func() {
		app.Get("/products?category_id=stickers").Do().Status(http.StatusBadRequest)
	})
})

var _ = Describe("GET /products?page=", func() {
	BeforeEach(func() {
		for range 5 {
			factories.Product(GinkgoT(), app.Store)
		}
	})

	DescribeTable("pages through the catalog",
		func(page int, wantProducts int) {
			var body views.ProductsResponse
			app.Get(fmt.Sprintf("/products?page=%d&per_page=2", page)).Do().
				Status(http.StatusOK).
				JSON(&body)
			Expect(body.Products).To(HaveLen(wantProducts))
			Expect(body.Total).To(Equal(int64(5)))
		},
		Entry("the first page is full", 1, 2),
		Entry("the last page holds the rest", 3, 1),
		Entry("a page past the end is empty", 4, 0),
	)

	It("refuses a per_page over the limit", func() {
		app.Get("/products?per_page=1000").Do().Status(http.StatusBadRequest)
	})
})
