package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Order.Validate", func() {
	var order models.Order

	BeforeEach(func() {
		order = factories.NewOrder()
	})

	Context("with the factory's defaults", func() {
		It("is valid", func() {
			Expect(order.Validate()).To(Succeed())
		})
	})

	Context("with a status other than placed", func() {
		BeforeEach(func() {
			order.Status = "shipped"
		})

		It("is refused", func() {
			Expect(order.Validate()).To(MatchError(models.ErrInvalid))
		})
	})

	Context("with totals that don't add up", func() {
		BeforeEach(func() {
			order.TotalCents = 1
		})

		It("is refused", func() {
			Expect(order.Validate()).To(MatchError(models.ErrInvalid))
		})
	})
})
