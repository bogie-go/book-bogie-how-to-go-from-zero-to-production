package adapter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Coupons", func() {
	Describe("GetCoupon", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				_, err := store.GetCoupon(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no coupon has that id", "999999999"),
		)
	})

	Describe("DeleteCoupon", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				Expect(store.DeleteCoupon(GinkgoT().Context(), id)).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no coupon has that id", "999999999"),
		)
	})

	Describe("UpdateCoupon", func() {
		Context("when the coupon is not there", func() {
			var missing models.Coupon

			BeforeEach(func() {
				missing = factories.NewCoupon()
				missing.ID = 999999999
			})

			It("is not found", func() {
				_, err := store.UpdateCoupon(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})

	Describe("CreateCoupon", func() {
		Context("when a reference names nothing", func() {
			var missing models.Coupon

			BeforeEach(func() {
				missing = factories.NewCoupon()
				missing.PriceRuleID = 999999999
			})

			It("is invalid", func() {
				_, err := store.CreateCoupon(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrInvalid))
			})
		})
	})

	// No round trip here: a coupon belongs to a parent that has to be saved
	// first. Make the parent with its factory and add the specs, with the
	// reference set in the factory's "with" function.
})
