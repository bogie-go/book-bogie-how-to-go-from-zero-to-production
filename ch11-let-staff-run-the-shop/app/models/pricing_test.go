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
			Expect(models.PercentageDiscount{Percentage: percentage}.Apply(ctx).CurrentCents).To(Equal(want))
		},
		Entry("20% of $50", int32(20), int32(5000), int32(4000)),
		Entry("of the price so far, not the original $100", int32(10), int32(5000), int32(4500)),
		Entry("from nothing", int32(20), int32(0), int32(0)),
		Entry("rounding the discount down to whole cents", int32(15), int32(333), int32(284)),
	)
})

var _ = Describe("BuyXGetY", func() {
	// The cart holds keycaps (product 5) and four stickers (product 1).
	DescribeTable("gives free keycaps, and leaves the price alone",
		func(keycaps int32, wantFree []models.FreeItem) {
			ctx := models.PriceContext{
				Lines: []models.CartLine{
					{ProductID: 5, Name: "Keycap", PriceCents: 1800, Quantity: keycaps},
					{ProductID: 1, Name: "Sticker", PriceCents: 300, Quantity: 4},
				},
				CurrentCents: 20000,
			}
			rule := models.BuyXGetY{ProductID: 5, Buy: 2, Free: 1}
			after := rule.Apply(ctx)
			Expect(after.FreeItems).To(Equal(wantFree))
			Expect(after.CurrentCents).To(Equal(int32(20000)))
		},
		Entry("one keycap: none free yet", int32(1), []models.FreeItem(nil)),
		Entry("two keycaps: one free", int32(2), []models.FreeItem{{ProductID: 5, Name: "Keycap", Quantity: 1}}),
		Entry("three keycaps: still one free", int32(3), []models.FreeItem{{ProductID: 5, Name: "Keycap", Quantity: 1}}),
		Entry("four keycaps: two free", int32(4), []models.FreeItem{{ProductID: 5, Name: "Keycap", Quantity: 2}}),
	)
})

var _ = Describe("PriceCart", func() {
	rule := func(name, kind, properties string) models.PriceRule {
		return models.PriceRule{Name: name, Kind: kind, Properties: []byte(properties)}
	}
	twentyOff := rule("20% off", "percentage_discount", `{"percentage":20}`)
	tenOff := rule("10% off", "percentage_discount", `{"percentage":10}`)
	keycapFree := rule("Keycaps", "buy_x_get_y", `{"product_id":5,"buy":2,"free":1}`)
	keycaps := []models.CartLine{{ProductID: 5, Name: "Keycap", PriceCents: 1800, Quantity: 3}}

	It("passes the price from rule to rule, and records each step", func() {
		priced := models.PriceCart(keycaps, []models.PriceRule{twentyOff, tenOff})
		Expect(priced.Applied).To(Equal([]models.AppliedRule{
			{Name: "20% off", BeforeCents: 5400, AfterCents: 4320},
			{Name: "10% off", BeforeCents: 4320, AfterCents: 3888},
		}))
		Expect(priced.TotalCents).To(Equal(int32(3888)))
		Expect(priced.DiscountCents()).To(Equal(int32(1512)))
	})

	It("gives free items, named after their rule, at no cost", func() {
		priced := models.PriceCart(keycaps, []models.PriceRule{twentyOff, keycapFree})
		Expect(priced.FreeItems).To(Equal([]models.FreeItem{{ProductID: 5, Name: "Keycap", Quantity: 1, Rule: "Keycaps"}}))
		Expect(priced.Applied).To(Equal([]models.AppliedRule{{Name: "20% off", BeforeCents: 5400, AfterCents: 4320}}))
	})

	It("gives the same price whichever runs first, a free item or a discount", func() {
		first := models.PriceCart(keycaps, []models.PriceRule{twentyOff, keycapFree})
		second := models.PriceCart(keycaps, []models.PriceRule{keycapFree, twentyOff})
		Expect(first.TotalCents).To(Equal(int32(4320)))
		Expect(second.TotalCents).To(Equal(first.TotalCents))
	})

	It("records only the rules that did something", func() {
		stickers := []models.CartLine{{ProductID: 1, PriceCents: 300, Quantity: 2}}
		priced := models.PriceCart(stickers, []models.PriceRule{twentyOff, keycapFree})
		Expect(priced.Applied).To(Equal([]models.AppliedRule{{Name: "20% off", BeforeCents: 600, AfterCents: 480}}))
		Expect(priced.FreeItems).To(BeEmpty())
	})

	It("skips a rule that won't load, and applies the rest", func() {
		broken := rule("Broken", "percentage_discount", `{"percentage":"ten"}`)
		priced := models.PriceCart(keycaps, []models.PriceRule{broken, twentyOff})
		Expect(priced.Skipped).To(HaveLen(1))
		Expect(priced.TotalCents).To(Equal(int32(4320)))
	})

	It("gives a different price when two discounts run in another order", func() {
		fiveOff := rule("$5 off $40", "minimum_purchase", `{"minimum_cents":4000,"discount_cents":500}`)
		cart := []models.CartLine{{ProductID: 2, PriceCents: 4500, Quantity: 1}}
		Expect(models.PriceCart(cart, []models.PriceRule{twentyOff, fiveOff}).TotalCents).To(Equal(int32(3600)))
		Expect(models.PriceCart(cart, []models.PriceRule{fiveOff, twentyOff}).TotalCents).To(Equal(int32(3200)))
	})

	It("never takes the price below zero", func() {
		fiveOffAnything := rule("$5 off", "minimum_purchase", `{"minimum_cents":0,"discount_cents":500}`)
		sticker := []models.CartLine{{ProductID: 1, PriceCents: 300, Quantity: 1}}
		Expect(models.PriceCart(sticker, []models.PriceRule{fiveOffAnything}).TotalCents).To(BeZero())
	})
})

var _ = Describe("MinimumPurchase", func() {
	DescribeTable("takes $5 off once the price so far reaches $40",
		func(original, current, want int32) {
			ctx := models.PriceContext{OriginalCents: original, CurrentCents: current}
			Expect(models.MinimumPurchase{MinimumCents: 4000, DiscountCents: 500}.Apply(ctx).CurrentCents).To(Equal(want))
		},
		Entry("under the minimum", int32(3999), int32(3999), int32(3999)),
		Entry("exactly the minimum", int32(4000), int32(4000), int32(3500)),
		Entry("over the minimum", int32(6000), int32(6000), int32(5500)),
		Entry("a $45 cart already brought down to $36", int32(4500), int32(3600), int32(3600)),
	)
})
