package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Product.Validate", func() {
	var product models.Product

	BeforeEach(func() {
		product = factories.NewProduct()
	})

	Context("with the factory's defaults", func() {
		It("is valid", func() {
			Expect(product.Validate()).To(Succeed())
		})
	})

	Context("with no name", func() {
		BeforeEach(func() { product.Name = "" })

		It("is refused", func() {
			err := product.Validate()
			Expect(err).To(MatchError(models.ErrInvalid))
			Expect(err).To(MatchError(ContainSubstring("name is required")))
		})
	})

	Context("with a price of zero", func() {
		BeforeEach(func() { product.PriceCents = 0 })

		It("is refused", func() {
			Expect(product.Validate()).To(MatchError(ContainSubstring("price must be more than zero")))
		})
	})

	Context("with negative stock", func() {
		BeforeEach(func() { product.Stock = -1 })

		It("is refused", func() {
			Expect(product.Validate()).To(MatchError(ContainSubstring("stock can't be negative")))
		})
	})

	Context("with no stock left", func() {
		BeforeEach(func() { product.Stock = 0 })

		It("is still valid", func() {
			Expect(product.Validate()).To(Succeed())
		})
	})
})
