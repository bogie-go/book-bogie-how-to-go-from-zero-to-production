package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("OrderItem.Validate", func() {
	var orderItem models.OrderItem

	BeforeEach(func() {
		orderItem = factories.NewOrderItem()
	})

	Context("with the factory's defaults", func() {
		It("is valid", func() {
			Expect(orderItem.Validate()).To(Succeed())
		})
	})

	Context("with a blank name", func() {
		BeforeEach(func() {
			orderItem.Name = ""
		})

		It("is refused", func() {
			Expect(orderItem.Validate()).To(MatchError(models.ErrInvalid))
		})
	})

	Context("with a quantity of zero", func() {
		BeforeEach(func() {
			orderItem.Quantity = 0
		})

		It("is refused", func() {
			Expect(orderItem.Validate()).To(MatchError(models.ErrInvalid))
		})
	})
})
