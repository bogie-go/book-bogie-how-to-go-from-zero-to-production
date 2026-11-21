package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
	"depot/test/testdb"
)

var _ = Describe("GET /cart", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodGet),
			HaveField("Path", "/cart"),
		)))
	})

	It("applies the sales running now, and shows what they give away", func() {
		keycap := factories.Product(GinkgoT(), app.Store, func(p *models.Product) {
			p.Name = "Bogie Keycap"
			p.ImageUrl = "/images/products/bogie-keycap.webp"
			p.PriceCents = 1800
			p.Stock = 10
		})
		factories.PriceRule(GinkgoT(), app.Store, func(r *models.PriceRule) {
			r.Name = "Keycaps: buy 2, get 1 free"
			r.Kind = "buy_x_get_y"
			r.Properties = json.RawMessage(fmt.Sprintf(`{"product_id":%d,"buy":2,"free":1}`, keycap.ID))
		})
		cart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), "guest:1")
		Expect(err).NotTo(HaveOccurred())
		Expect(app.Store.AddCartItem(GinkgoT().Context(), cart.ID, keycap.ID, 2)).To(Succeed())

		app.Get("/cart").As("guest:1").Do().
			Status(http.StatusOK).
			Contains(`"subtotal_cents":3600,"discount_cents":0,"total_cents":3600,"applied_rules":[],` +
				fmt.Sprintf(`"free_items":[{"product_id":%d,"name":"Bogie Keycap",`, keycap.ID) +
				`"image_url":"/images/products/bogie-keycap.webp","quantity":1,"rule":"Keycaps: buy 2, get 1 free"}]`)
	})

	It("still answers when a rule in the database won't load", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.PriceCents = 1000 })
		broken := factories.PriceRule(GinkgoT(), app.Store)
		_, err := testdb.Pool(GinkgoT()).Exec(GinkgoT().Context(),
			`UPDATE price_rules SET properties = '{"percentage":"ten"}' WHERE id = $1`, broken.ID)
		Expect(err).NotTo(HaveOccurred())
		factories.PriceRule(GinkgoT(), app.Store, func(r *models.PriceRule) { r.Name = "10% off" })

		cart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), "guest:1")
		Expect(err).NotTo(HaveOccurred())
		Expect(app.Store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 1)).To(Succeed())

		app.Get("/cart").As("guest:1").Do().
			Status(http.StatusOK).
			Contains(`"applied_rules":[{"name":"10% off","price_before_cents":1000,"price_after_cents":900,"discount_cents":100}]`)
	})
})
