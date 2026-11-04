package adapter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("OrderItems", func() {
	Describe("GetOrderItem", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				_, err := store.GetOrderItem(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no order item has that id", "999999999"),
		)
	})

	Describe("DeleteOrderItem", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				Expect(store.DeleteOrderItem(GinkgoT().Context(), id)).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no order item has that id", "999999999"),
		)
	})

	Describe("UpdateOrderItem", func() {
		Context("when the order item is not there", func() {
			var missing models.OrderItem

			BeforeEach(func() {
				missing = factories.NewOrderItem()
				missing.ID = 999999999
			})

			It("is not found", func() {
				_, err := store.UpdateOrderItem(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})

	Describe("CreateOrderItem", func() {
		Context("when a reference names nothing", func() {
			var missing models.OrderItem

			BeforeEach(func() {
				missing = factories.NewOrderItem()
				missing.OrderID = 999999999
				missing.ProductID = 999999999
			})

			It("is invalid", func() {
				_, err := store.CreateOrderItem(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrInvalid))
			})
		})
	})

	// No round trip here: an order item belongs to a parent that has to be saved
	// first. Make the parent with its factory and add the specs, with the
	// reference set in the factory's "with" function.
})
