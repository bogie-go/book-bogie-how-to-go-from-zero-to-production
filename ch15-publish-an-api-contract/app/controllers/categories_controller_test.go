package controllers_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("GET /categories", func() {
	Context("with no categories", func() {
		It("is an empty list", func() {
			app.Get("/categories").Do().Status(http.StatusOK).Contains("\"categories\":[]")
		})
	})

	It("counts the products customers can buy, and lists every category", func() {
		stickers := factories.Category(GinkgoT(), app.Store, func(c *models.Category) { c.Name = "Stickers" })
		desk := factories.Category(GinkgoT(), app.Store, func(c *models.Category) { c.Name = "Desk" })
		factories.Category(GinkgoT(), app.Store, func(c *models.Category) { c.Name = "Drinkware" })

		inCategories := func(archived bool, categoryIDs ...int64) {
			p := factories.Product(GinkgoT(), app.Store)
			Expect(app.Store.SetProductCategories(GinkgoT().Context(), fmt.Sprint(p.ID), categoryIDs)).To(Succeed())
			if archived {
				_, err := app.Store.ArchiveProduct(GinkgoT().Context(), fmt.Sprint(p.ID))
				Expect(err).NotTo(HaveOccurred())
			}
		}
		inCategories(false, stickers.ID)
		inCategories(true, stickers.ID)
		inCategories(true, desk.ID)

		app.Get("/categories").Do().
			Status(http.StatusOK).
			Contains(`"name":"Desk","product_count":0`).
			Contains(`"name":"Drinkware","product_count":0`).
			Contains(`"name":"Stickers","product_count":1`)
	})
})
