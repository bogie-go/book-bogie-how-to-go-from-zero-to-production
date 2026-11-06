package controllers_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("POST /users", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodPost),
			HaveField("Path", "/users"),
		)))
	})

	It("creates an account, already signed in", func() {
		app.Post("/users", map[string]any{"email": "Ada@Example.com", "password": "correct-horse"}).Do().
			Status(http.StatusCreated).
			Contains(`"email":"ada@example.com"`).
			Contains(`"token":"`)
	})

	It("refuses an email that's already taken", func() {
		factories.User(GinkgoT(), app.Store, func(u *models.User) { u.Email = "ada@example.com" })

		app.Post("/users", map[string]any{"email": "ada@example.com", "password": "correct-horse"}).Do().
			Status(http.StatusUnprocessableEntity).
			Contains(`"error":"invalid: email is already taken"`)
	})

	It("refuses a password that's too short", func() {
		app.Post("/users", map[string]any{"email": "ada@example.com", "password": "short"}).Do().
			Status(http.StatusUnprocessableEntity).
			Contains("password must be 8 to 72 bytes long")
	})

	It("brings the guest's cart along", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Stock = 10 })
		guestCart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), "guest:1")
		Expect(err).NotTo(HaveOccurred())
		Expect(app.Store.AddCartItem(GinkgoT().Context(), guestCart.ID, mug.ID, 2)).To(Succeed())

		app.Post("/users", map[string]any{"email": "grace@example.com", "password": "correct-horse"}).
			As("guest:1").Do().
			Status(http.StatusCreated)

		grace, err := app.Store.FindUserByEmail(GinkgoT().Context(), "grace@example.com")
		Expect(err).NotTo(HaveOccurred())
		app.Get("/cart").As(fmt.Sprintf("user:%d", grace.ID)).Do().Contains(`"quantity":2`)
	})
})
