package order_confirmation_test

import (
	"context"
	"io"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/app/services/mailer"
	"depot/app/services/order_confirmation"
	"depot/test/factories"
	"depot/test/testdb"
)

// sentMail is a sender that records every message instead of sending it.
type sentMail struct{ messages []mailer.Message }

func (s *sentMail) Send(_ context.Context, msg mailer.Message) error {
	s.messages = append(s.messages, msg)
	return nil
}

var _ = Describe("Send", func() {
	It("emails the customer their receipt", func() {
		store := testdb.Store(GinkgoT())
		ada := factories.User(GinkgoT(), store, func(u *models.User) { u.Email = "ada@example.com" })
		order := factories.Order(GinkgoT(), store, func(o *models.Order) {
			o.UserID = ada.ID
			o.SubtotalCents, o.DiscountCents, o.TotalCents = 2800, 0, 2800
		})
		mug := factories.Product(GinkgoT(), store)
		factories.OrderItem(GinkgoT(), store, func(i *models.OrderItem) {
			i.OrderID, i.ProductID = order.ID, mug.ID
			i.Name, i.Quantity, i.UnitPriceCents = "Bogie Ceramic Mug", 2, 1400
		})
		mail := &sentMail{}
		service := order_confirmation.New(store, mail, slog.New(slog.NewTextHandler(io.Discard, nil)))

		Expect(service.Send(GinkgoT().Context(), order.ID)).To(Succeed())

		Expect(mail.messages).To(HaveLen(1))
		Expect(mail.messages[0].To).To(Equal("ada@example.com"))
		Expect(mail.messages[0].Body).To(ContainSubstring("2 x Bogie Ceramic Mug  $28.00"))
		Expect(mail.messages[0].Body).To(ContainSubstring("Total     $28.00"))
	})
})
