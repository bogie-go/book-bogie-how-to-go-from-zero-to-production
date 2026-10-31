package password_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/app/services/password"
)

var _ = Describe("Hash", func() {
	DescribeTable("refuses a password of the wrong length",
		func(plain string) {
			_, err := password.Hash(plain)
			Expect(err).To(MatchError(models.ErrInvalid))
		},
		Entry("7 bytes", "1234567"),
		Entry("73 bytes", strings.Repeat("a", 73)),
		Entry("37 accented letters, which are 74 bytes", strings.Repeat("é", 37)),
	)

	It("makes a digest that Check accepts, and only for that password", func() {
		digest, err := password.Hash("correct-horse")
		Expect(err).NotTo(HaveOccurred())
		Expect(digest).NotTo(ContainSubstring("correct-horse"))
		Expect(password.Check(digest, "correct-horse")).To(BeTrue())
		Expect(password.Check(digest, "wrong-horse")).To(BeFalse())
	})

	It("makes a different digest every time, thanks to the salt", func() {
		first, _ := password.Hash("correct-horse")
		second, _ := password.Hash("correct-horse")
		Expect(first).NotTo(Equal(second))
	})
})
