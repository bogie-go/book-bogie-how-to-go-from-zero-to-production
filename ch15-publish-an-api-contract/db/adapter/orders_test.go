package adapter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Orders", func() {
	Describe("GetOrder", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				_, err := store.GetOrder(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no order has that id", "999999999"),
		)
	})

	Describe("DeleteOrder", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				Expect(store.DeleteOrder(GinkgoT().Context(), id)).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no order has that id", "999999999"),
		)
	})

	Describe("UpdateOrder", func() {
		Context("when the order is not there", func() {
			var missing models.Order

			BeforeEach(func() {
				missing = factories.NewOrder()
				missing.ID = 999999999
			})

			It("is not found", func() {
				_, err := store.UpdateOrder(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})

	Describe("CreateOrder", func() {
		Context("when a reference names nothing", func() {
			var missing models.Order

			BeforeEach(func() {
				missing = factories.NewOrder()
				missing.UserID = 999999999
			})

			It("is invalid", func() {
				_, err := store.CreateOrder(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrInvalid))
			})
		})
	})

	// No round trip here: an order belongs to a parent that has to be saved
	// first. Make the parent with its factory and add the specs, with the
	// reference set in the factory's "with" function.
})
