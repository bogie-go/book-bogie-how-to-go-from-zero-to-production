package controllers_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
	"depot/test/testdb"
)

var _ = Describe("GET /orders/:id", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodGet),
			HaveField("Path", "/orders/:id"),
		)))
	})

	It("shows another customer's order as not found", func() {
		owner := factories.User(GinkgoT(), app.Store)
		order := factories.Order(GinkgoT(), app.Store, func(o *models.Order) { o.UserID = owner.ID })
		stranger := factories.User(GinkgoT(), app.Store)

		app.Get(fmt.Sprintf("/orders/%d", order.ID)).As(fmt.Sprintf("user:%d", stranger.ID)).Do().
			Status(http.StatusNotFound)
		app.Get(fmt.Sprintf("/orders/%d", order.ID)).As(fmt.Sprintf("user:%d", owner.ID)).Do().
			Status(http.StatusOK)
	})
})

var _ = Describe("POST /orders", func() {
	var ada string
	var mug models.Product

	BeforeEach(func() {
		user := factories.User(GinkgoT(), app.Store)
		ada = fmt.Sprintf("user:%d", user.ID)
		mug = factories.Product(GinkgoT(), app.Store, func(p *models.Product) {
			p.Name = "Bogie Ceramic Mug"
			p.PriceCents = 1400
			p.Stock = 3
		})
		cart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), ada)
		Expect(err).NotTo(HaveOccurred())
		Expect(app.Store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 2)).To(Succeed())
	})

	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodPost),
			HaveField("Path", "/orders"),
		)))
	})

	It("places an order from the cart", func() {
		app.Post("/orders", map[string]any{}).As(ada).Do().
			Status(http.StatusCreated).
			Contains(`"status":"placed"`).
			Contains(`"subtotal_cents":2800,"discount_cents":0,"total_cents":2800`)

		app.Get("/cart").As(ada).Do().Contains(`"items":[]`)
		app.Get(fmt.Sprintf("/products/%d", mug.ID)).Do().Contains(`"stock":1`)
	})

	It("refuses an empty cart", func() {
		app.Post("/orders", map[string]any{}).As(ada).Do().Status(http.StatusCreated)
		app.Post("/orders", map[string]any{}).As(ada).Do().
			Status(http.StatusUnprocessableEntity).
			Contains("the cart is empty")
	})

	It("refuses a guest", func() {
		app.Post("/orders", map[string]any{}).As("guest:1").Do().
			Status(http.StatusForbidden).
			Contains("create an account to place an order")
	})

	It("changes nothing when an item runs out halfway", func() {
		keycap := factories.Product(GinkgoT(), app.Store, func(p *models.Product) {
			p.Name = "Bogie Keycap"
			p.Stock = 1
		})
		cart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), ada)
		Expect(err).NotTo(HaveOccurred())
		Expect(app.Store.AddCartItem(GinkgoT().Context(), cart.ID, keycap.ID, 1)).To(Succeed())
		_, err = testdb.Pool(GinkgoT()).Exec(GinkgoT().Context(),
			`UPDATE products SET stock = 0 WHERE id = $1`, keycap.ID)
		Expect(err).NotTo(HaveOccurred())

		app.Post("/orders", map[string]any{}).As(ada).Do().
			Status(http.StatusUnprocessableEntity).
			Contains("not enough Bogie Keycap left")

		app.Get(fmt.Sprintf("/products/%d", mug.ID)).Do().Contains(`"stock":3`)
		app.Get("/cart").As(ada).Do().Contains(`"quantity":2`)
	})

	It("applies a coupon code, in any case", func() {
		rule := factories.PriceRule(GinkgoT(), app.Store, func(r *models.PriceRule) {
			r.Kind = "amount_off"
			r.Properties = []byte(`{"amount_cents":500,"min_subtotal_cents":0}`)
		})
		factories.Coupon(GinkgoT(), app.Store, func(c *models.Coupon) {
			c.Code = "LAUNCH5"
			c.PriceRuleID = rule.ID
		})

		app.Post("/orders", map[string]any{"coupon_code": "launch5"}).As(ada).Do().
			Status(http.StatusCreated).
			Contains(`"subtotal_cents":2800,"discount_cents":500,"total_cents":2300`)
	})

	It("refuses a code that isn't valid, and keeps the cart", func() {
		app.Post("/orders", map[string]any{"coupon_code": "NOPE"}).As(ada).Do().
			Status(http.StatusUnprocessableEntity)
		app.Get("/cart").As(ada).Do().Contains(`"quantity":2`)
	})

	It("refuses the code of a coupon whose rule has ended", func() {
		rule := factories.PriceRule(GinkgoT(), app.Store, func(r *models.PriceRule) {
			r.StartsAt = time.Now().Add(-48 * time.Hour)
			r.EndsAt = time.Now().Add(-24 * time.Hour)
		})
		factories.Coupon(GinkgoT(), app.Store, func(c *models.Coupon) {
			c.Code = "EXPIRED"
			c.PriceRuleID = rule.ID
		})

		app.Post("/orders", map[string]any{"coupon_code": "EXPIRED"}).As(ada).Do().
			Status(http.StatusUnprocessableEntity)
	})
})
