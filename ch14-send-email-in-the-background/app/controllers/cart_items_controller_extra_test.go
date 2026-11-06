package controllers_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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

		app.Post("/cart_items", body).As("guest:1").Do().Status(http.StatusOK)
		app.Post("/cart_items", body).As("guest:1").Do().
			Status(http.StatusOK).
			Contains(`"quantity":2`).
			Contains(`"total_cents":2800`)
	})
})

var _ = Describe("Someone else's cart item", func() {
	var itemPath string

	BeforeEach(func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Stock = 10 })
		cart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), "guest:alice")
		Expect(err).NotTo(HaveOccurred())
		Expect(app.Store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 1)).To(Succeed())
		lines, err := app.Store.CartContents(GinkgoT().Context(), cart.ID)
		Expect(err).NotTo(HaveOccurred())
		itemPath = fmt.Sprintf("/cart_items/%d", lines[0].ItemID)
	})

	It("can't be changed by another visitor", func() {
		app.Patch(itemPath, map[string]any{"quantity": 5}).As("guest:mallory").Do().
			Status(http.StatusNotFound)
	})

	It("can't be removed by another visitor", func() {
		app.Delete(itemPath).As("guest:mallory").Do().
			Status(http.StatusNotFound)
	})
})
