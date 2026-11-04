package adapter_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("PriceRules", func() {
	Describe("GetPriceRule", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				_, err := store.GetPriceRule(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no price rule has that id", "999999999"),
		)

		Context("with a saved price rule", func() {
			var record models.PriceRule

			BeforeEach(func() {
				record = factories.PriceRule(GinkgoT(), store)
			})

			It("finds the price rule by id", func() {
				Expect(store.GetPriceRule(GinkgoT().Context(), strconv.FormatInt(record.ID, 10))).To(Equal(record))
			})
		})
	})

	Describe("ListPriceRules", func() {
		Context("with a saved price rule", func() {
			var record models.PriceRule

			BeforeEach(func() {
				record = factories.PriceRule(GinkgoT(), store)
			})

			It("includes the saved price rule", func() {
				all, err := store.ListPriceRules(GinkgoT().Context(), 10)
				Expect(err).NotTo(HaveOccurred())
				Expect(all).To(HaveLen(1))
				Expect(all[0].ID).To(Equal(record.ID))
			})
		})
	})

	Describe("UpdatePriceRule", func() {
		Context("with a saved price rule", func() {
			var record models.PriceRule

			BeforeEach(func() {
				record = factories.PriceRule(GinkgoT(), store)
			})

			It("stores the change", func() {
				record.Name = "Edited"
				updated, err := store.UpdatePriceRule(GinkgoT().Context(), record)
				Expect(err).NotTo(HaveOccurred())
				Expect(updated.ID).To(Equal(record.ID))
				Expect(updated.Name).To(BeEquivalentTo("Edited"))
			})
		})

		Context("when the price rule is not there", func() {
			var missing models.PriceRule

			BeforeEach(func() {
				missing = factories.NewPriceRule()
				missing.ID = 999999999
			})

			It("is not found", func() {
				_, err := store.UpdatePriceRule(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})

	Describe("DeletePriceRule", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				Expect(store.DeletePriceRule(GinkgoT().Context(), id)).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no price rule has that id", "999999999"),
		)

		Context("with a saved price rule", func() {
			var record models.PriceRule

			BeforeEach(func() {
				record = factories.PriceRule(GinkgoT(), store)
			})

			It("removes the price rule", func() {
				id := strconv.FormatInt(record.ID, 10)
				Expect(store.DeletePriceRule(GinkgoT().Context(), id)).To(Succeed())
				_, err := store.GetPriceRule(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})
})
