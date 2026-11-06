package adapter_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Users", func() {
	Describe("GetUser", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				_, err := store.GetUser(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no user has that id", "999999999"),
		)

		Context("with a saved user", func() {
			var record models.User

			BeforeEach(func() {
				record = factories.User(GinkgoT(), store)
			})

			It("finds the user by id", func() {
				Expect(store.GetUser(GinkgoT().Context(), strconv.FormatInt(record.ID, 10))).To(Equal(record))
			})
		})
	})

	Describe("ListUsers", func() {
		Context("with a saved user", func() {
			var record models.User

			BeforeEach(func() {
				record = factories.User(GinkgoT(), store)
			})

			It("includes the saved user", func() {
				all, err := store.ListUsers(GinkgoT().Context(), 10)
				Expect(err).NotTo(HaveOccurred())
				Expect(all).To(HaveLen(1))
				Expect(all[0].ID).To(Equal(record.ID))
			})
		})
	})

	Describe("UpdateUser", func() {
		Context("with a saved user", func() {
			var record models.User

			BeforeEach(func() {
				record = factories.User(GinkgoT(), store)
			})

			It("stores the change", func() {
				record.Email = "edited@example.com"
				updated, err := store.UpdateUser(GinkgoT().Context(), record)
				Expect(err).NotTo(HaveOccurred())
				Expect(updated.ID).To(Equal(record.ID))
				Expect(updated.Email).To(BeEquivalentTo("edited@example.com"))
			})
		})

		Context("when the user is not there", func() {
			var missing models.User

			BeforeEach(func() {
				missing = factories.NewUser()
				missing.ID = 999999999
			})

			It("is not found", func() {
				_, err := store.UpdateUser(GinkgoT().Context(), missing)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})

	Describe("DeleteUser", func() {
		DescribeTable("returns ErrNotFound",
			func(id string) {
				Expect(store.DeleteUser(GinkgoT().Context(), id)).To(MatchError(models.ErrNotFound))
			},
			Entry("when the id is not a number", "nope"),
			Entry("when the id is empty", ""),
			Entry("when no user has that id", "999999999"),
		)

		Context("with a saved user", func() {
			var record models.User

			BeforeEach(func() {
				record = factories.User(GinkgoT(), store)
			})

			It("removes the user", func() {
				id := strconv.FormatInt(record.ID, 10)
				Expect(store.DeleteUser(GinkgoT().Context(), id)).To(Succeed())
				_, err := store.GetUser(GinkgoT().Context(), id)
				Expect(err).To(MatchError(models.ErrNotFound))
			})
		})
	})
})
