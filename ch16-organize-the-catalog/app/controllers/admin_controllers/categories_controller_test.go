package admin_controllers_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/app/views"
	"depot/test/factories"
)

var _ = Describe("GET /admin/categories", func() {
	Context("with no categories", func() {
		It("is an empty list", func() {
			app.Get("/admin/categories").As(staff).Do().Status(http.StatusOK).Contains("\"categories\":[]")
		})
	})

	Context("with categories", func() {
		var first, second models.Category

		BeforeEach(func() {
			first = factories.Category(GinkgoT(), app.Store)
			second = factories.Category(GinkgoT(), app.Store)
		})

		It("lists them, newest first", func() {
			var list views.AdminCategoriesResponse
			app.Get("/admin/categories").As(staff).Do().Status(http.StatusOK).JSON(&list)

			Expect(list.Categories).To(HaveLen(2))
			Expect(list.Categories[0].ID).To(Equal(second.ID))
			Expect(list.Categories[1].ID).To(Equal(first.ID))
		})
	})
})

var _ = Describe("GET /admin/categories/:id", func() {
	Context("when the category exists", func() {
		var record models.Category

		BeforeEach(func() {
			record = factories.Category(GinkgoT(), app.Store)
		})

		It("shows it", func() {
			var shown views.AdminCategoryResponse
			app.Get("/admin/categories/" + fmt.Sprint(record.ID)).As(staff).Do().Status(http.StatusOK).JSON(&shown)

			Expect(shown.ID).To(Equal(record.ID))
		})
	})

	Context("when the category does not exist", func() {
		It("is not found", func() {
			app.Get("/admin/categories/999999999").As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

var _ = Describe("POST /admin/categories", func() {
	var body views.AdminCategoryRequest

	Context("with a valid body", func() {
		BeforeEach(func() {
			body = categoryRequest(factories.NewCategory())
		})

		It("creates it", func() {
			var created views.AdminCategoryResponse
			app.Post("/admin/categories", body).As(staff).Do().Status(http.StatusCreated).JSON(&created)

			Expect(created.ID).NotTo(BeZero())
			_, err := app.Store.GetCategory(GinkgoT().Context(), fmt.Sprint(created.ID))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("with a malformed body", func() {
		It("is a bad request", func() {
			app.Post("/admin/categories", "{\"nope\":").As(staff).Do().Status(http.StatusBadRequest)
		})
	})
})

var _ = Describe("PATCH /admin/categories/:id", func() {
	Context("when the category exists", func() {
		var record models.Category

		BeforeEach(func() {
			record = factories.Category(GinkgoT(), app.Store)
		})

		It("changes only the fields sent", func() {
			var updated views.AdminCategoryResponse
			app.Patch("/admin/categories/"+fmt.Sprint(record.ID), map[string]any{"name": "Edited"}).As(staff).Do().Status(http.StatusOK).JSON(&updated)

			Expect(updated.Name).To(BeEquivalentTo("Edited"))
		})

		Context("with a malformed body", func() {
			It("is a bad request", func() {
				app.Patch("/admin/categories/"+fmt.Sprint(record.ID), "{\"nope\":").As(staff).Do().Status(http.StatusBadRequest)
			})
		})
	})

	Context("when the category does not exist", func() {
		It("is not found", func() {
			app.Patch("/admin/categories/999999999", map[string]any{}).As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

var _ = Describe("DELETE /admin/categories/:id", func() {
	Context("when the category exists", func() {
		var record models.Category

		BeforeEach(func() {
			record = factories.Category(GinkgoT(), app.Store)
		})

		It("deletes it", func() {
			app.Delete("/admin/categories/" + fmt.Sprint(record.ID)).As(staff).Do().Status(http.StatusNoContent)

			app.Get("/admin/categories/" + fmt.Sprint(record.ID)).As(staff).Do().Status(http.StatusNotFound)
		})
	})

	Context("when the category does not exist", func() {
		It("is not found", func() {
			app.Delete("/admin/categories/999999999").As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

// categoryRequest is the body that creates a record like this one.
func categoryRequest(record models.Category) views.AdminCategoryRequest {
	return views.AdminCategoryRequest{Name: record.Name}
}
