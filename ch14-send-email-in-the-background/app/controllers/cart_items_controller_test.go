package controllers_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /cart_items", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodPost),
			HaveField("Path", "/cart_items"),
		)))
	})
})

var _ = Describe("PATCH /cart_items/:id", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodPatch),
			HaveField("Path", "/cart_items/:id"),
		)))
	})
})

var _ = Describe("DELETE /cart_items/:id", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodDelete),
			HaveField("Path", "/cart_items/:id"),
		)))
	})
})
