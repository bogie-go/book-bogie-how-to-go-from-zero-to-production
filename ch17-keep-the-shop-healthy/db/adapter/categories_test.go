package adapter_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Categories", func() {
	Describe("GetCategory", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				_, err := store.GetCategory(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no category has that id", "999999999"),
		)

		Context("with a saved category", func() {
			var record models.Category

			BeforeEach(func() {
				record = factories.Category(GinkgoT(), store)
			})

			It("finds the category by id", func() {
				Expect(store.GetCategory(GinkgoT().Context(), strconv.FormatInt(record.ID, 10))).To(Equal(record))
			})
		})
	})

	Describe("ListCategories", func() {
		Context("with a saved category", func() {
			var record models.Category

			BeforeEach(func() {
				record = factories.Category(GinkgoT(), store)
			})

			It("includes the saved category", func() {
				all, err := store.ListCategories(GinkgoT().Context(), 10)
				Expect(err).NotTo(HaveOccurred())
				Expect(all).To(HaveLen(1))
				Expect(all[0].ID).To(Equal(record.ID))
			})
		})
	})

	Describe("UpdateCategory", func() {
		Context("with a saved category", func() {
			var record models.Category

			BeforeEach(func() {
				record = factories.Category(GinkgoT(), store)
			})

			It("stores the change", func() {
				record.Name = "Edited"
				updated, err := store.UpdateCategory(GinkgoT().Context(), record)
				Expect(err).NotTo(HaveOccurred())
				Expect(updated.ID).To(Equal(record.ID))
				Expect(updated.Name).To(BeEquivalentTo("Edited"))
			})
		})

		Context("when the category is not there", func() {
			var missing models.Category

			BeforeEach(func() {
				missing = factories.NewCategory()
				missing.ID = 999999999
			})

			It("is not found", func() {
				_, err := store.UpdateCategory(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})

	Describe("DeleteCategory", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				Expect(store.DeleteCategory(GinkgoT().Context(), id)).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no category has that id", "999999999"),
		)

		Context("with a saved category", func() {
			var record models.Category

			BeforeEach(func() {
				record = factories.Category(GinkgoT(), store)
			})

			It("removes the category", func() {
				id := strconv.FormatInt(record.ID, 10)
				Expect(store.DeleteCategory(GinkgoT().Context(), id)).To(Succeed())
				_, err := store.GetCategory(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})
})
