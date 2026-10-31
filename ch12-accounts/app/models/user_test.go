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
})
