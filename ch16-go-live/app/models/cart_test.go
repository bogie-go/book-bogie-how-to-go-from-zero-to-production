package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Cart.Validate", func() {
	var cart models.Cart

	BeforeEach(func() {
		cart = factories.NewCart()
	})

	Context("with the factory's defaults", func() {
		It("is valid", func() {
			Expect(cart.Validate()).To(Succeed())
		})
	})

	// Add a Context for each rule you put in Validate, for example a
	// required token:
	//
	//	Context("with a blank token", func() {
	//		BeforeEach(func() {
	//			cart.Token = ""
	//		})
	//
	//		It("is refused", func() {
	//			Expect(cart.Validate()).To(MatchError(models.ErrInvalid))
	//		})
	//	})
})
