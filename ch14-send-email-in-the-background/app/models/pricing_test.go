package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
)

var _ = Describe("PercentOff", func() {
	DescribeTable("takes a percentage off the whole cart",
		func(percent int32, lines []models.CartLine, want int32) {
			Expect(models.PercentOff{Percent: percent}.Discount(lines)).To(Equal(want))
		},
		Entry("20% of $50", int32(20), []models.CartLine{{PriceCents: 2500, Quantity: 2}}, int32(1000)),
		Entry("from an empty cart", int32(20), []models.CartLine{}, int32(0)),
		Entry("rounding down to whole cents", int32(15), []models.CartLine{{PriceCents: 333, Quantity: 1}}, int32(49)),
	)
})

var _ = Describe("BuyXGetY", func() {
	DescribeTable("makes some keycaps free, and nothing else",
		func(keycaps int32, want int32) {
			lines := []models.CartLine{
				{ProductID: 5, PriceCents: 1800, Quantity: keycaps},
				{ProductID: 1, PriceCents: 300, Quantity: 4},
			}
			rule := models.BuyXGetY{ProductID: 5, Buy: 2, Free: 1}
			Expect(rule.Discount(lines)).To(Equal(want))
		},
		Entry("two keycaps: none free yet", int32(2), int32(0)),
		Entry("three keycaps: one free", int32(3), int32(1800)),
		Entry("six keycaps: two free", int32(6), int32(3600)),
		Entry("seven keycaps: still two free", int32(7), int32(3600)),
	)
})

var _ = Describe("AmountOff", func() {
	DescribeTable("takes $5 off orders of $40 or more",
		func(subtotalCents int32, want int32) {
			lines := []models.CartLine{{PriceCents: subtotalCents, Quantity: 1}}
			Expect(models.AmountOff{AmountCents: 500, MinSubtotalCents: 4000}.Discount(lines)).To(Equal(want))
		},
		Entry("under the minimum", int32(3999), int32(0)),
		Entry("exactly the minimum", int32(4000), int32(500)),
		Entry("over the minimum", int32(6000), int32(500)),
	)
})

var _ = Describe("PriceCart", func() {
	rule := func(name, kind, properties string) models.PriceRule {
		return models.PriceRule{Name: name, Kind: kind, Properties: []byte(properties)}
	}
	lines := []models.CartLine{{ProductID: 1, PriceCents: 300, Quantity: 2}}

	It("lists only the rules that took something off", func() {
		priced := models.PriceCart(lines, []models.PriceRule{
			rule("10% off", "percent_off", `{"percent":10}`),
			rule("Keycaps", "buy_x_get_y", `{"product_id":5,"buy":2,"free":1}`),
		})
		Expect(priced.Applied).To(Equal([]models.AppliedRule{{Name: "10% off", DiscountCents: 60}}))
		Expect(priced.TotalCents).To(Equal(int32(540)))
	})

	It("never discounts more than the subtotal", func() {
		priced := models.PriceCart(lines, []models.PriceRule{
			rule("90% off", "percent_off", `{"percent":90}`),
			rule("50% off", "percent_off", `{"percent":50}`),
		})
		Expect(priced.DiscountCents).To(Equal(int32(600)))
		Expect(priced.TotalCents).To(BeZero())
	})

	It("skips a rule that won't load, and applies the rest", func() {
		priced := models.PriceCart(lines, []models.PriceRule{
			rule("Broken", "percent_off", `{"percent":"ten"}`),
			rule("10% off", "percent_off", `{"percent":10}`),
		})
		Expect(priced.Skipped).To(HaveLen(1))
		Expect(priced.Applied).To(Equal([]models.AppliedRule{{Name: "10% off", DiscountCents: 60}}))
	})
})
