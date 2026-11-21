package order_confirmation_test

import (
	"io"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/app/services/order_confirmation"
	"depot/test/factories"
	"depot/test/testdb"
)

// Not printed in the book: Chapter 14 says a line with a unit price of zero,
// one a rule gave away, reads "free" on the receipt. This checks it.
var _ = Describe("Send, with a free item", func() {
	It("prints the free line as free, not $0.00", func() {
		store := testdb.Store(GinkgoT())
		ada := factories.User(GinkgoT(), store, func(u *models.User) { u.Email = "ada@example.com" })
		order := factories.Order(GinkgoT(), store, func(o *models.Order) {
			o.UserID = ada.ID
			o.SubtotalCents, o.DiscountCents, o.TotalCents = 3600, 0, 3600
		})
		keycap := factories.Product(GinkgoT(), store)
		factories.OrderItem(GinkgoT(), store, func(i *models.OrderItem) {
			i.OrderID, i.ProductID = order.ID, keycap.ID
			i.Name, i.Quantity, i.UnitPriceCents = "Bogie Keycap", 2, 1800
		})
		factories.OrderItem(GinkgoT(), store, func(i *models.OrderItem) {
			i.OrderID, i.ProductID = order.ID, keycap.ID
			i.Name, i.Quantity, i.UnitPriceCents = "Bogie Keycap", 1, 0
		})
		mail := &sentMail{}
		service := order_confirmation.New(store, mail, slog.New(slog.NewTextHandler(io.Discard, nil)))

		Expect(service.Send(GinkgoT().Context(), order.ID)).To(Succeed())

		Expect(mail.messages).To(HaveLen(1))
		Expect(mail.messages[0].Body).To(ContainSubstring("2 x Bogie Keycap  $36.00\n1 x Bogie Keycap  free\n"))
		Expect(mail.messages[0].Body).NotTo(ContainSubstring("$0.00"))
		Expect(mail.messages[0].Body).To(ContainSubstring("Total     $36.00"))
	})
})
