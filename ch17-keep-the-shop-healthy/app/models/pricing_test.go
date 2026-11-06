package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
)

var _ = Describe("PercentageDiscount", func() {
	DescribeTable("takes a percentage off the price so far",
		func(percentage, current, want int32) {
			ctx := models.PriceContext{OriginalCents: 10000, CurrentCents: current}
			Expect(models.PercentageDiscount{Percentage: percentage}.Apply(ctx)).To(Equal(want))
		},
		Entry("20% of $50", int32(20), int32(5000), int32(4000)),
		Entry("of the price so far, not the original $100", int32(10), int32(5000), int32(4500)),
		Entry("from nothing", int32(20), int32(0), int32(0)),
		Entry("rounding the discount down to whole cents", int32(15), int32(333), int32(284)),
	)
})

var _ = Describe("BuyXGetY", func() {
	// The cart holds keycaps (product 5) and four stickers (product 1).
	DescribeTable("takes the free keycaps off, and nothing else",
		func(keycaps, wantOff int32) {
			ctx := models.PriceContext{
				Lines: []models.CartLine{
					{ProductID: 5, PriceCents: 1800, Quantity: keycaps},
					{ProductID: 1, PriceCents: 300, Quantity: 4},
				},
				CurrentCents: 20000,
			}
			rule := models.BuyXGetY{ProductID: 5, Buy: 2, Free: 1}
			Expect(ctx.CurrentCents - rule.Apply(ctx)).To(Equal(wantOff))
		},
		Entry("two keycaps: none free yet", int32(2), int32(0)),
		Entry("three keycaps: one free", int32(3), int32(1800)),
		Entry("six keycaps: two free", int32(6), int32(3600)),
		Entry("seven keycaps: still two free", int32(7), int32(3600)),
	)
})

var _ = Describe("PriceCart", func() {
	rule := func(name, kind, properties string) models.PriceRule {
		return models.PriceRule{Name: name, Kind: kind, Properties: []byte(properties)}
	}
	twentyOff := rule("20% off", "percentage_discount", `{"percentage":20}`)
	keycapFree := rule("Keycaps", "buy_x_get_y", `{"product_id":5,"buy":2,"free":1}`)
	keycaps := []models.CartLine{{ProductID: 5, PriceCents: 1800, Quantity: 3}}

	It("passes the price from rule to rule, and records each step", func() {
		priced := models.PriceCart(keycaps, []models.PriceRule{twentyOff, keycapFree})
		Expect(priced.Applied).To(Equal([]models.AppliedRule{
			{Name: "20% off", BeforeCents: 5400, AfterCents: 4320},
			{Name: "Keycaps", BeforeCents: 4320, AfterCents: 2520},
		}))
		Expect(priced.TotalCents).To(Equal(int32(2520)))
		Expect(priced.DiscountCents()).To(Equal(int32(2880)))
	})

	It("gives a different price when the same rules run in another order", func() {
		priced := models.PriceCart(keycaps, []models.PriceRule{keycapFree, twentyOff})
		Expect(priced.TotalCents).To(Equal(int32(2880)))
	})

	It("records only the rules that changed the price", func() {
		stickers := []models.CartLine{{ProductID: 1, PriceCents: 300, Quantity: 2}}
		priced := models.PriceCart(stickers, []models.PriceRule{twentyOff, keycapFree})
		Expect(priced.Applied).To(Equal([]models.AppliedRule{{Name: "20% off", BeforeCents: 600, AfterCents: 480}}))
	})

	It("never takes the price below zero", func() {
		ninetyOff := rule("90% off", "percentage_discount", `{"percentage":90}`)
		priced := models.PriceCart(keycaps, []models.PriceRule{ninetyOff, keycapFree})
		Expect(priced.TotalCents).To(BeZero())
	})

	It("skips a rule that won't load, and applies the rest", func() {
		broken := rule("Broken", "percentage_discount", `{"percentage":"ten"}`)
		priced := models.PriceCart(keycaps, []models.PriceRule{broken, twentyOff})
		Expect(priced.Skipped).To(HaveLen(1))
		Expect(priced.TotalCents).To(Equal(int32(4320)))
	})
})

var _ = Describe("MinimumPurchase", func() {
	DescribeTable("takes $5 off once the price so far reaches $40",
		func(original, current, want int32) {
			ctx := models.PriceContext{OriginalCents: original, CurrentCents: current}
			Expect(models.MinimumPurchase{MinimumCents: 4000, DiscountCents: 500}.Apply(ctx)).To(Equal(want))
		},
		Entry("under the minimum", int32(3999), int32(3999), int32(3999)),
		Entry("exactly the minimum", int32(4000), int32(4000), int32(3500)),
		Entry("over the minimum", int32(6000), int32(6000), int32(5500)),
		Entry("a $45 cart already brought down to $36", int32(4500), int32(3600), int32(3600)),
	)
})
