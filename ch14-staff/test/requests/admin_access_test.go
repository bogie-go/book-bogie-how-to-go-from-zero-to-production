package requests_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("The admin area", func() {
	It("refuses a request with no token", func() {
		app.Get("/admin/products").Do().Status(http.StatusUnauthorized)
	})

	It("refuses a customer", func() {
		customer := factories.User(GinkgoT(), app.Store)
		app.Get("/admin/products").As(fmt.Sprintf("user:%d", customer.ID)).Do().
			Status(http.StatusForbidden).
			Contains("staff only")
	})

	It("refuses a guest, even with a valid guest token", func() {
		app.Get("/admin/products").As("guest:1").Do().
			Status(http.StatusForbidden).
			Contains("staff only")
	})

	It("lets staff in", func() {
		staff := factories.User(GinkgoT(), app.Store, func(u *models.User) { u.Role = models.RoleStaff })
		app.Get("/admin/products").As(fmt.Sprintf("user:%d", staff.ID)).Do().
			Status(http.StatusOK)
	})

	It("refuses staff on the very next request after their role is taken away", func() {
		staff := factories.User(GinkgoT(), app.Store, func(u *models.User) { u.Role = models.RoleStaff })
		subject := fmt.Sprintf("user:%d", staff.ID)
		app.Get("/admin/products").As(subject).Do().Status(http.StatusOK)

		staff.Role = ""
		_, err := app.Store.UpdateUser(GinkgoT().Context(), staff)
		Expect(err).NotTo(HaveOccurred())

		app.Get("/admin/products").As(subject).Do().Status(http.StatusForbidden)
	})
})
