package jobs_test

import (
	"context"
	"io"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"depot/app/jobs"
)

// sentOrders is a confirmer that records the order ids it was asked to send.
type sentOrders struct{ ids []int64 }

func (s *sentOrders) Send(_ context.Context, orderID int64) error {
	s.ids = append(s.ids, orderID)
	return nil
}

var _ = Describe("SendOrderConfirmationWorker", func() {
	var (
		confirm *sentOrders
		worker  *jobs.SendOrderConfirmationWorker
		job     *river.Job[jobs.SendOrderConfirmationArgs]
	)

	Describe("Work", func() {
		Context("with a job for an order", func() {
			BeforeEach(func() {
				confirm = &sentOrders{}
				worker = &jobs.SendOrderConfirmationWorker{Confirm: confirm, Log: slog.New(slog.NewTextHandler(io.Discard, nil))}
				job = &river.Job[jobs.SendOrderConfirmationArgs]{JobRow: &rivertype.JobRow{ID: 1}, Args: jobs.SendOrderConfirmationArgs{OrderID: 42}}
			})

			It("sends that order's confirmation", func() {
				Expect(worker.Work(context.Background(), job)).To(Succeed())
				Expect(confirm.ids).To(Equal([]int64{42}))
			})
		})
	})
})
