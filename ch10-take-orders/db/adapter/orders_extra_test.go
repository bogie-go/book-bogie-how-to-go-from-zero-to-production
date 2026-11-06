package adapter_test

import (
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("PlaceOrder at the same moment", func() {
	// cartFor makes a customer whose cart holds the given products, one of
	// each, added in that order.
	cartFor := func(products ...models.Product) (userID, cartID int64) {
		user := factories.User(GinkgoT(), store)
		cart, err := store.FindOrCreateCartByToken(GinkgoT().Context(), fmt.Sprintf("user:%d", user.ID))
		Expect(err).NotTo(HaveOccurred())
		for _, p := range products {
			Expect(store.AddCartItem(GinkgoT().Context(), cart.ID, p.ID, 1)).To(Succeed())
		}
		return user.ID, cart.ID
	}

	It("places one order when the same cart is checked out many times at once", func() {
		mug := factories.Product(GinkgoT(), store, func(p *models.Product) { p.Stock = 1000 })
		for range 10 {
			userID, cartID := cartFor(mug)

			errs := make([]error, 8)
			var wg sync.WaitGroup
			for i := range errs {
				wg.Go(func() {
					_, errs[i] = store.PlaceOrder(GinkgoT().Context(), userID, cartID, nil)
				})
			}
			wg.Wait()

			placed := 0
			for _, err := range errs {
				if err == nil {
					placed++
				} else {
					Expect(err).To(MatchError(ContainSubstring("the cart is empty")))
				}
			}
			Expect(placed).To(Equal(1))
		}
	})

	It("sells the last item only once", func() {
		mug := factories.Product(GinkgoT(), store, func(p *models.Product) { p.Stock = 1 })
		errs := make([]error, 8)
		users, carts := make([]int64, 8), make([]int64, 8)
		for i := range errs {
			users[i], carts[i] = cartFor(mug)
		}

		var wg sync.WaitGroup
		for i := range errs {
			wg.Go(func() {
				_, errs[i] = store.PlaceOrder(GinkgoT().Context(), users[i], carts[i], nil)
			})
		}
		wg.Wait()

		placed := 0
		for _, err := range errs {
			if err == nil {
				placed++
			}
		}
		Expect(placed).To(Equal(1))
		left, err := store.GetProduct(GinkgoT().Context(), fmt.Sprint(mug.ID))
		Expect(err).NotTo(HaveOccurred())
		Expect(left.Stock).To(BeZero())
	})

	It("never deadlocks when two carts hold the same products in opposite order", func() {
		for range 30 {
			a := factories.Product(GinkgoT(), store, func(p *models.Product) { p.Stock = 10 })
			b := factories.Product(GinkgoT(), store, func(p *models.Product) { p.Stock = 10 })
			user1, cart1 := cartFor(a, b)
			user2, cart2 := cartFor(b, a)

			errs := make([]error, 2)
			var wg sync.WaitGroup
			wg.Go(func() { _, errs[0] = store.PlaceOrder(GinkgoT().Context(), user1, cart1, nil) })
			wg.Go(func() { _, errs[1] = store.PlaceOrder(GinkgoT().Context(), user2, cart2, nil) })
			wg.Wait()

			Expect(errs).To(HaveEach(BeNil()))
		}
	})
})
