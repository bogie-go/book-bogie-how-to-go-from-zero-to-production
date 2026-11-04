package admin_controllers_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("GET /admin/coupons", func() {
	Context("with no coupons", func() {
		It("is an empty list", func() {
			app.Get("/admin/coupons").As(staff).Do().Status(http.StatusOK).Contains("\"coupons\":[]")
		})
	})
})

var _ = Describe("GET /admin/coupons/:id", func() {
	Context("when the coupon does not exist", func() {
		It("is not found", func() {
			app.Get("/admin/coupons/999999999").As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

var _ = Describe("POST /admin/coupons", func() {
	Context("with a malformed body", func() {
		It("is a bad request", func() {
			app.Post("/admin/coupons", "{\"nope\":").As(staff).Do().Status(http.StatusBadRequest)
		})
	})
})

var _ = Describe("PATCH /admin/coupons/:id", func() {
	Context("when the coupon does not exist", func() {
		It("is not found", func() {
			app.Patch("/admin/coupons/999999999", map[string]any{}).As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

var _ = Describe("DELETE /admin/coupons/:id", func() {
	Context("when the coupon does not exist", func() {
		It("is not found", func() {
			app.Delete("/admin/coupons/999999999").As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

// The specs that save a coupon are left out: it belongs to a parent that has to
// be saved first. Make the parent with its factory, set the reference in the
// "with" function of factories.Coupon, and add them.

var _ = Describe("Staff saving a coupon", func() {
	It("can't point a coupon at a rule that doesn't exist", func() {
		app.Post("/admin/coupons", map[string]any{"code": "GHOST", "price_rule_id": 999999999}).As(staff).Do().
			Status(http.StatusUnprocessableEntity).
			Contains("price_rule_id does not exist")
	})
})
