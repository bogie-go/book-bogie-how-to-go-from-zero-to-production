package controllers_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/app/services/password"
	"depot/test/factories"
)

var _ = Describe("POST /sessions", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodPost),
			HaveField("Path", "/sessions"),
		)))
	})

	BeforeEach(func() {
		digest, err := password.Hash("correct-horse")
		Expect(err).NotTo(HaveOccurred())
		factories.User(GinkgoT(), app.Store, func(u *models.User) {
			u.Email = "ada@example.com"
			u.PasswordDigest = digest
		})
	})

	It("signs in with the right password", func() {
		app.Post("/sessions", map[string]any{"email": "ada@example.com", "password": "correct-horse"}).Do().
			Status(http.StatusOK).
			Contains(`"token":"`)
	})

	DescribeTable("refuses wrong credentials, the same way every time",
		func(email, plain string) {
			app.Post("/sessions", map[string]any{"email": email, "password": plain}).Do().
				Status(http.StatusUnauthorized).
				Contains(`"error":"email or password is incorrect"`)
		},
		Entry("a wrong password", "ada@example.com", "wrong-horse"),
		Entry("an email with no account", "bob@example.com", "correct-horse"),
	)

	It("brings the guest's cart along", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Stock = 10 })
		guestCart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), "guest:1")
		Expect(err).NotTo(HaveOccurred())
		Expect(app.Store.AddCartItem(GinkgoT().Context(), guestCart.ID, mug.ID, 2)).To(Succeed())

		app.Post("/sessions", map[string]any{"email": "ada@example.com", "password": "correct-horse"}).
			As("guest:1").Do().
			Status(http.StatusOK)

		ada, err := app.Store.FindUserByEmail(GinkgoT().Context(), "ada@example.com")
		Expect(err).NotTo(HaveOccurred())
		app.Get("/cart").As(fmt.Sprintf("user:%d", ada.ID)).Do().Contains(`"quantity":2`)
		app.Get("/cart").As("guest:1").Do().Contains(`"items":[]`)
	})

	It("never merges another customer's cart", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Stock = 10 })
		theirCart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), "user:99")
		Expect(err).NotTo(HaveOccurred())
		Expect(app.Store.AddCartItem(GinkgoT().Context(), theirCart.ID, mug.ID, 3)).To(Succeed())

		app.Post("/sessions", map[string]any{"email": "ada@example.com", "password": "correct-horse"}).
			As("user:99").Do().
			Status(http.StatusOK)

		app.Get("/cart").As("user:99").Do().Contains(`"quantity":3`)
	})

	It("signs in even with a garbage token", func() {
		app.Post("/sessions", map[string]any{"email": "ada@example.com", "password": "correct-horse"}).
			WithAuthorization("Bearer not-a-token").Do().
			Status(http.StatusOK)
	})
})
