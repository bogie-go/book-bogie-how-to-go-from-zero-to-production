package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

var _ = Describe("Category.Validate", func() {
	var category models.Category

	BeforeEach(func() {
		category = factories.NewCategory()
	})

	Context("with the factory's defaults", func() {
		It("is valid", func() {
			Expect(category.Validate()).To(Succeed())
		})
	})

	Context("with a blank name", func() {
		BeforeEach(func() { category.Name = "" })

		It("is refused", func() {
			Expect(category.Validate()).To(MatchError(models.ErrInvalid))
		})
	})
})
