package adapter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("AddCartItem", func() {
	var cart models.Cart
	var mug models.Product

	BeforeEach(func() {
		cart = factories.Cart(GinkgoT(), store)
		mug = factories.Product(GinkgoT(), store, func(p *models.Product) { p.Stock = 3 })
	})

	It("adds a new line", func() {
		Expect(store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 1)).To(Succeed())
		lines, err := store.CartContents(GinkgoT().Context(), cart.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(lines).To(HaveLen(1))
		Expect(lines[0].Quantity).To(Equal(int32(1)))
	})

	It("adds to the existing line for the same product", func() {
		Expect(store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 1)).To(Succeed())
		Expect(store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 2)).To(Succeed())
		lines, err := store.CartContents(GinkgoT().Context(), cart.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(lines).To(HaveLen(1))
		Expect(lines[0].Quantity).To(Equal(int32(3)))
	})

	It("refuses more than the product has in stock", func() {
		Expect(store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 2)).To(Succeed())
		err := store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 2)
		Expect(err).To(MatchError(models.ErrInvalid))
		Expect(err).To(MatchError(ContainSubstring("only 3 left in stock")))
	})

	It("refuses a product that doesn't exist", func() {
		Expect(store.AddCartItem(GinkgoT().Context(), cart.ID, 999999999, 1)).To(MatchError(models.ErrNotFound))
	})
})
