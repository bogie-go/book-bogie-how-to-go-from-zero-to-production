package controllers_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/test/testdb"
)

var _ = Describe("GET /status", func() {
	It("is mounted", func() {
		Expect(app.Server.Router.Routes()).To(ContainElement(And(
			HaveField("Method", http.MethodGet),
			HaveField("Path", "/status"),
		)))
	})

	It("says ok when the database answers and no job needs attention", func() {
		app.Get("/status").Do().
			Status(http.StatusOK).
			Contains(`"status":"ok","database":"ok","jobs_waiting_too_long":0,"jobs_failed_today":0`)
	})

	It("asks for attention when a job has failed for good", func() {
		_, err := testdb.Pool(GinkgoT()).Exec(GinkgoT().Context(), `
			INSERT INTO river_job (kind, args, state, finalized_at, max_attempts)
			VALUES ('send_order_confirmation', '{"order_id": 1}', 'discarded', now(), 25)`)
		Expect(err).NotTo(HaveOccurred())

		app.Get("/status").Do().
			Status(http.StatusServiceUnavailable).
			Contains(`"jobs_failed_today":1`)
	})
})
