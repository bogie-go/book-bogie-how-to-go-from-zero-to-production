package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("CartItem.Validate", func() {
	var cartItem models.CartItem

	BeforeEach(func() {
		cartItem = factories.NewCartItem()
	})

	Context("with the factory's defaults", func() {
		It("is valid", func() {
			Expect(cartItem.Validate()).To(Succeed())
		})
	})

	Context("with a quantity of zero", func() {
		BeforeEach(func() { cartItem.Quantity = 0 })

		It("is refused", func() {
			Expect(cartItem.Validate()).To(MatchError(ContainSubstring("quantity must be at least 1")))
		})
	})
})
