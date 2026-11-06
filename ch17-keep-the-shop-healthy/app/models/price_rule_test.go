package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("PriceRule.Validate", func() {
	var priceRule models.PriceRule

	BeforeEach(func() {
		priceRule = factories.NewPriceRule()
	})

	Context("with the factory's defaults", func() {
		It("is valid", func() {
			Expect(priceRule.Validate()).To(Succeed())
		})
	})

	DescribeTable("refuses settings that don't fit the kind",
		func(kind, properties string) {
			priceRule.Kind = kind
			priceRule.Properties = []byte(properties)
			Expect(priceRule.Validate()).To(MatchError(models.ErrInvalid))
		},
		Entry("an unknown kind", "half_price", `{}`),
		Entry("no settings at all", "percentage_discount", `{}`),
		Entry("a misspelled setting beside a good one", "percentage_discount", `{"percentage":10,"percntage":10}`),
		Entry("a setting of the wrong type", "percentage_discount", `{"percentage":"ten"}`),
		Entry("more than 100 percent", "percentage_discount", `{"percentage":150}`),
		Entry("a negative percentage", "percentage_discount", `{"percentage":-10}`),
		Entry("buy nothing, get one free", "buy_x_get_y", `{"product_id":5,"buy":0,"free":1}`),
		Entry("a negative discount", "minimum_purchase", `{"minimum_cents":4000,"discount_cents":-500}`),
	)
})
