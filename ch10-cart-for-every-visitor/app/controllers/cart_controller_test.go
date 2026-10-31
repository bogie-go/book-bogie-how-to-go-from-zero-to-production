package controllers_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GET /cart", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodGet),
			HaveField("Path", "/cart"),
		)))
	})
})
