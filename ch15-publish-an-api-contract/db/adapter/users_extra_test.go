package adapter_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("A user's role", func() {
	It("is saved and read back", func() {
		saved := factories.User(GinkgoT(), store, func(u *models.User) { u.Role = models.RoleStaff })

		found, err := store.GetUser(GinkgoT().Context(), fmt.Sprint(saved.ID))
		Expect(err).NotTo(HaveOccurred())
		Expect(found.Role).To(Equal(models.RoleStaff))
	})
})
