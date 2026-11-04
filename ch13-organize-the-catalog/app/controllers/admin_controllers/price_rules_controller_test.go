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

var _ = Describe("GET /admin/price_rules", func() {
	Context("with no price_rules", func() {
		It("is an empty list", func() {
			app.Get("/admin/price_rules").As(staff).Do().Status(http.StatusOK).Contains("\"price_rules\":[]")
		})
	})

	Context("with price_rules", func() {
		var first, second models.PriceRule

		BeforeEach(func() {
			first = factories.PriceRule(GinkgoT(), app.Store)
			second = factories.PriceRule(GinkgoT(), app.Store)
		})

		It("lists them, newest first", func() {
			var list views.AdminPriceRulesResponse
			app.Get("/admin/price_rules").As(staff).Do().Status(http.StatusOK).JSON(&list)

			Expect(list.PriceRules).To(HaveLen(2))
			Expect(list.PriceRules[0].ID).To(Equal(second.ID))
			Expect(list.PriceRules[1].ID).To(Equal(first.ID))
		})
	})
})

var _ = Describe("GET /admin/price_rules/:id", func() {
	Context("when the price_rule exists", func() {
		var record models.PriceRule

		BeforeEach(func() {
			record = factories.PriceRule(GinkgoT(), app.Store)
		})

		It("shows it", func() {
			var shown views.AdminPriceRuleResponse
			app.Get("/admin/price_rules/" + fmt.Sprint(record.ID)).As(staff).Do().Status(http.StatusOK).JSON(&shown)

			Expect(shown.ID).To(Equal(record.ID))
		})
	})

	Context("when the price_rule does not exist", func() {
		It("is not found", func() {
			app.Get("/admin/price_rules/999999999").As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

var _ = Describe("POST /admin/price_rules", func() {
	var body views.AdminPriceRuleRequest

	Context("with a valid body", func() {
		BeforeEach(func() {
			body = priceRuleRequest(factories.NewPriceRule())
		})

		It("creates it", func() {
			var created views.AdminPriceRuleResponse
			app.Post("/admin/price_rules", body).As(staff).Do().Status(http.StatusCreated).JSON(&created)

			Expect(created.ID).NotTo(BeZero())
			_, err := app.Store.GetPriceRule(GinkgoT().Context(), fmt.Sprint(created.ID))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("with a malformed body", func() {
		It("is a bad request", func() {
			app.Post("/admin/price_rules", "{\"nope\":").As(staff).Do().Status(http.StatusBadRequest)
		})
	})
})

var _ = Describe("PATCH /admin/price_rules/:id", func() {
	Context("when the price_rule exists", func() {
		var record models.PriceRule

		BeforeEach(func() {
			record = factories.PriceRule(GinkgoT(), app.Store)
		})

		It("changes only the fields sent", func() {
			var updated views.AdminPriceRuleResponse
			app.Patch("/admin/price_rules/"+fmt.Sprint(record.ID), map[string]any{"name": "Edited"}).As(staff).Do().Status(http.StatusOK).JSON(&updated)

			Expect(updated.Name).To(BeEquivalentTo("Edited"))
			Expect(updated.Kind).To(Equal(record.Kind))
			Expect(updated.Priority).To(Equal(record.Priority))
		})

		Context("with a malformed body", func() {
			It("is a bad request", func() {
				app.Patch("/admin/price_rules/"+fmt.Sprint(record.ID), "{\"nope\":").As(staff).Do().Status(http.StatusBadRequest)
			})
		})
	})

	Context("when the price_rule does not exist", func() {
		It("is not found", func() {
			app.Patch("/admin/price_rules/999999999", map[string]any{}).As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

var _ = Describe("DELETE /admin/price_rules/:id", func() {
	Context("when the price_rule exists", func() {
		var record models.PriceRule

		BeforeEach(func() {
			record = factories.PriceRule(GinkgoT(), app.Store)
		})

		It("deletes it", func() {
			app.Delete("/admin/price_rules/" + fmt.Sprint(record.ID)).As(staff).Do().Status(http.StatusNoContent)

			app.Get("/admin/price_rules/" + fmt.Sprint(record.ID)).As(staff).Do().Status(http.StatusNotFound)
		})
	})

	Context("when the price_rule does not exist", func() {
		It("is not found", func() {
			app.Delete("/admin/price_rules/999999999").As(staff).Do().Status(http.StatusNotFound)
		})
	})
})

// priceRuleRequest is the body that creates a record like this one.
func priceRuleRequest(record models.PriceRule) views.AdminPriceRuleRequest {
	return views.AdminPriceRuleRequest{Name: record.Name, Kind: record.Kind, Properties: record.Properties, Priority: record.Priority, StartsAt: record.StartsAt, EndsAt: record.EndsAt}
}

var _ = Describe("Staff saving a price rule", func() {
	It("is held to the same rules as everything else", func() {
		app.Post("/admin/price_rules", map[string]any{
			"name": "Too generous", "kind": "percent_off", "properties": map[string]any{"percent": 150},
			"priority": 1, "starts_at": "2026-11-01T00:00:00Z", "ends_at": "2026-12-01T00:00:00Z",
		}).As(staff).Do().
			Status(http.StatusUnprocessableEntity).
			Contains("percent must be between 1 and 100")
	})
})
