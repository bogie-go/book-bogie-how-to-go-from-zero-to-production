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
		Entry("no settings at all", "percent_off", `{}`),
		Entry("a misspelled setting beside a good one", "percent_off", `{"percent":10,"percnt":10}`),
		Entry("a setting of the wrong type", "percent_off", `{"percent":"ten"}`),
		Entry("more than 100 percent", "percent_off", `{"percent":150}`),
		Entry("a negative percent", "percent_off", `{"percent":-10}`),
		Entry("buy nothing, get one free", "buy_x_get_y", `{"product_id":5,"buy":0,"free":1}`),
		Entry("a negative amount off", "amount_off", `{"amount_cents":-500,"min_subtotal_cents":4000}`),
	)
})
