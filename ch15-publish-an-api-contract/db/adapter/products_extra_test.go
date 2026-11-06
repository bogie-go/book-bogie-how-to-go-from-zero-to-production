package adapter_test

import (
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("EditProduct during checkouts", func() {
	It("never undoes the stock a checkout took", func() {
		for range 5 {
			mug := factories.Product(GinkgoT(), store, func(p *models.Product) { p.Stock = 100 })
			carts := make([][2]int64, 8)
			for i := range carts {
				user := factories.User(GinkgoT(), store)
				cart, err := store.FindOrCreateCartByToken(GinkgoT().Context(), fmt.Sprintf("user:%d", user.ID))
				Expect(err).NotTo(HaveOccurred())
				Expect(store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 1)).To(Succeed())
				carts[i] = [2]int64{user.ID, cart.ID}
			}

			var wg sync.WaitGroup
			for _, uc := range carts {
				wg.Go(func() { _, _ = store.PlaceOrder(GinkgoT().Context(), uc[0], uc[1], nil, nil) })
				wg.Go(func() {
					_, _ = store.EditProduct(GinkgoT().Context(), fmt.Sprint(mug.ID), func(p *models.Product) error {
						p.PriceCents++
						return nil
					})
				})
			}
			wg.Wait()

			after, err := store.GetProduct(GinkgoT().Context(), fmt.Sprint(mug.ID))
			Expect(err).NotTo(HaveOccurred())
			Expect(after.Stock).To(Equal(int32(92)))
		}
	})
})
