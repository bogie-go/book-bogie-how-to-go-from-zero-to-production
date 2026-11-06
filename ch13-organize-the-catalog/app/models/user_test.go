package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("User.Validate", func() {
	var user models.User

	BeforeEach(func() {
		user = factories.NewUser()
	})

	Context("with the factory's defaults", func() {
		It("is valid", func() {
			Expect(user.Validate()).To(Succeed())
		})
	})

	Context("with an email that has no @", func() {
		BeforeEach(func() {
			user.Email = "ada.example.com"
		})

		It("is refused", func() {
			Expect(user.Validate()).To(MatchError(models.ErrInvalid))
		})
	})

	Context("with no password digest", func() {
		BeforeEach(func() {
			user.PasswordDigest = ""
		})

		It("is refused", func() {
			Expect(user.Validate()).To(MatchError(models.ErrInvalid))
		})
	})

	Context("with a role that isn't one of the shop's", func() {
		BeforeEach(func() { user.Role = "Staff" })

		It("is refused", func() {
			Expect(user.Validate()).To(MatchError(models.ErrInvalid))
		})
	})
})

var _ = Describe("UserIDFromSubject", func() {
	DescribeTable("reads a customer's ID, and nothing else",
		func(subject string, wantID int64, wantOK bool) {
			id, ok := models.UserIDFromSubject(subject)
			Expect(ok).To(Equal(wantOK))
			Expect(id).To(Equal(wantID))
		},
		Entry("a customer", "user:42", int64(42), true),
		Entry("a guest", "guest:3f2b8c1e", int64(0), false),
		Entry("a malformed customer", "user:forty-two", int64(0), false),
		Entry("an empty subject", "", int64(0), false),
	)

	It("reads back what UserSubject writes", func() {
		id, ok := models.UserIDFromSubject(models.UserSubject(7))
		Expect(ok).To(BeTrue())
		Expect(id).To(Equal(int64(7)))
	})
})
