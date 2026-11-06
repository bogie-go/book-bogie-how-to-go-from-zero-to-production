package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Coupon.Validate", func() {
	var coupon models.Coupon

	BeforeEach(func() {
		coupon = factories.NewCoupon()
	})

	Context("with the factory's defaults", func() {
		It("is valid", func() {
			Expect(coupon.Validate()).To(Succeed())
		})
	})

	Context("with a code in lower case", func() {
		BeforeEach(func() { coupon.Code = "launch5" })

		It("is refused", func() {
			Expect(coupon.Validate()).To(MatchError(models.ErrInvalid))
		})
	})
})
