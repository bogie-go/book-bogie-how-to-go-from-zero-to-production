package adapter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("CartItems", func() {
	Describe("GetCartItem", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				_, err := store.GetCartItem(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no cart item has that id", "999999999"),
		)
	})

	Describe("DeleteCartItem", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				Expect(store.DeleteCartItem(GinkgoT().Context(), id)).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no cart item has that id", "999999999"),
		)
	})

	Describe("UpdateCartItem", func() {
		Context("when the cart item is not there", func() {
			var missing models.CartItem

			BeforeEach(func() {
				missing = factories.NewCartItem()
				missing.ID = 999999999
			})

			It("is not found", func() {
				_, err := store.UpdateCartItem(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})

	Describe("CreateCartItem", func() {
		Context("when a reference names nothing", func() {
			var missing models.CartItem

			BeforeEach(func() {
				missing = factories.NewCartItem()
				missing.CartID = 999999999
				missing.ProductID = 999999999
			})

			It("is invalid", func() {
				_, err := store.CreateCartItem(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrInvalid))
			})
		})
	})

	// No round trip here: a cart item belongs to a parent that has to be saved
	// first. Make the parent with its factory and add the specs, with the
	// reference set in the factory's "with" function.
})
