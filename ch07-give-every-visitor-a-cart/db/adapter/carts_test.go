package adapter_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Carts", func() {
	Describe("GetCart", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				_, err := store.GetCart(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no cart has that id", "999999999"),
		)

		Context("with a saved cart", func() {
			var record models.Cart

			BeforeEach(func() {
				record = factories.Cart(GinkgoT(), store)
			})

			It("finds the cart by id", func() {
				Expect(store.GetCart(GinkgoT().Context(), strconv.FormatInt(record.ID, 10))).To(Equal(record))
			})
		})
	})

	Describe("ListCarts", func() {
		Context("with a saved cart", func() {
			var record models.Cart

			BeforeEach(func() {
				record = factories.Cart(GinkgoT(), store)
			})

			It("includes the saved cart", func() {
				all, err := store.ListCarts(GinkgoT().Context(), 10)
				Expect(err).NotTo(HaveOccurred())
				Expect(all).To(HaveLen(1))
				Expect(all[0].ID).To(Equal(record.ID))
			})
		})
	})

	Describe("UpdateCart", func() {
		Context("with a saved cart", func() {
			var record models.Cart

			BeforeEach(func() {
				record = factories.Cart(GinkgoT(), store)
			})

			It("stores the change", func() {
				record.Token = "Edited"
				updated, err := store.UpdateCart(GinkgoT().Context(), record)
				Expect(err).NotTo(HaveOccurred())
				Expect(updated.ID).To(Equal(record.ID))
				Expect(updated.Token).To(BeEquivalentTo("Edited"))
			})
		})

		Context("when the cart is not there", func() {
			var missing models.Cart

			BeforeEach(func() {
				missing = factories.NewCart()
				missing.ID = 999999999
			})

			It("is not found", func() {
				_, err := store.UpdateCart(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})

	Describe("DeleteCart", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				Expect(store.DeleteCart(GinkgoT().Context(), id)).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no cart has that id", "999999999"),
		)

		Context("with a saved cart", func() {
			var record models.Cart

			BeforeEach(func() {
				record = factories.Cart(GinkgoT(), store)
			})

			It("removes the cart", func() {
				id := strconv.FormatInt(record.ID, 10)
				Expect(store.DeleteCart(GinkgoT().Context(), id)).To(Succeed())
				_, err := store.GetCart(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})
})
