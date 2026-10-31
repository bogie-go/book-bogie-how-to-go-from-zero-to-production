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

var _ = Describe("GET /admin/products", func() {
	Context("with no products", func() {
		It("is an empty list", func() {
			app.Get("/admin/products").As(staff).Do().Status(http.StatusOK).Contains("\"products\":[]")
		})
	})

	Context("with products", func() {
		var first, second models.Product

		BeforeEach(func() {
			first = factories.Product(GinkgoT(), app.Store)
			second = factories.Product(GinkgoT(), app.Store)
		})

		It("lists them, newest first", func() {
			var list views.AdminProductsResponse
			app.Get("/admin/products").As(staff).Do().Status(http.StatusOK).JSON(&list)

			Expect(list.Products).To(HaveLen(2))
			Expect(list.Products[0].ID).To(Equal(second.ID))
			Expect(list.Products[1].ID).To(Equal(first.ID))
		})
	})
})

var _ = Describe("GET /admin/products/:id", func() {
	Context("when the product exists", func() {
		var record models.Product

		BeforeEach(func() {
			record = factories.Product(GinkgoT(), app.Store)
		})

		It("shows it", func() {
			var shown views.AdminProductResponse
			app.Get("/admin/products/" + fmt.Sprint(record.ID)).As(staff).Do().Status(http.StatusOK).JSON(&shown)

			Expect(shown.ID).To(Equal(record.ID))
		})
	})

	Context("when the product does not exist", func() {
		It("is not found", func() {
			app.Get("/admin/products/999999999").As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

var _ = Describe("POST /admin/products", func() {
	var body views.AdminProductRequest

	Context("with a valid body", func() {
		BeforeEach(func() {
			body = productRequest(factories.NewProduct())
		})

		It("creates it", func() {
			var created views.AdminProductResponse
			app.Post("/admin/products", body).As(staff).Do().Status(http.StatusCreated).JSON(&created)

			Expect(created.ID).NotTo(BeZero())
			_, err := app.Store.GetProduct(GinkgoT().Context(), fmt.Sprint(created.ID))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("with a malformed body", func() {
		It("is a bad request", func() {
			app.Post("/admin/products", "{\"nope\":").As(staff).Do().Status(http.StatusBadRequest)
		})
	})
})

var _ = Describe("PATCH /admin/products/:id", func() {
	Context("when the product exists", func() {
		var record models.Product

		BeforeEach(func() {
			record = factories.Product(GinkgoT(), app.Store)
		})

		It("changes only the fields sent", func() {
			var updated views.AdminProductResponse
			app.Patch("/admin/products/"+fmt.Sprint(record.ID), map[string]any{"name": "Edited"}).As(staff).Do().Status(http.StatusOK).JSON(&updated)

			Expect(updated.Name).To(BeEquivalentTo("Edited"))
			Expect(updated.Description).To(Equal(record.Description))
			Expect(updated.PriceCents).To(Equal(record.PriceCents))
			Expect(updated.Stock).To(Equal(record.Stock))
			Expect(updated.ImageUrl).To(Equal(record.ImageUrl))
		})

		Context("with a malformed body", func() {
			It("is a bad request", func() {
				app.Patch("/admin/products/"+fmt.Sprint(record.ID), "{\"nope\":").As(staff).Do().Status(http.StatusBadRequest)
			})
		})
	})

	Context("when the product does not exist", func() {
		It("is not found", func() {
			app.Patch("/admin/products/999999999", map[string]any{}).As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

var _ = Describe("DELETE /admin/products/:id", func() {
	Context("when the product exists", func() {
		var record models.Product

		BeforeEach(func() {
			record = factories.Product(GinkgoT(), app.Store)
		})

		It("deletes it", func() {
			app.Delete("/admin/products/" + fmt.Sprint(record.ID)).As(staff).Do().Status(http.StatusNoContent)

			app.Get("/admin/products/" + fmt.Sprint(record.ID)).As(staff).Do().Status(http.StatusNotFound)
		})
	})

	Context("when the product does not exist", func() {
		It("is not found", func() {
			app.Delete("/admin/products/999999999").As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

// productRequest is the body that creates a record like this one.
func productRequest(record models.Product) views.AdminProductRequest {
	return views.AdminProductRequest{Name: record.Name, Description: record.Description, PriceCents: record.PriceCents, Stock: record.Stock, ImageUrl: record.ImageUrl}
}

var _ = Describe("POST /admin/products/:id/archive and /restore", func() {
	It("archives a product, then brings it back", func() {
		mug := factories.Product(GinkgoT(), app.Store)
		path := fmt.Sprintf("/admin/products/%d", mug.ID)

		app.Post(path+"/archive", nil).As(staff).Do().
			Status(http.StatusOK).
			Contains(`"archived_at":"`)
		app.Post(path+"/restore", nil).As(staff).Do().
			Status(http.StatusOK).
			Contains(`"archived_at":null`)
	})
})
