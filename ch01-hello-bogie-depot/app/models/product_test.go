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

	// Add a Context for each rule you put in Validate, for example a
	// required name:
	//
	//	Context("with a blank name", func() {
	//		BeforeEach(func() {
	//			product.Name = ""
	//		})
	//
	//		It("is refused", func() {
	//			Expect(product.Validate()).To(MatchError(models.ErrInvalid))
	//		})
	//	})
})
