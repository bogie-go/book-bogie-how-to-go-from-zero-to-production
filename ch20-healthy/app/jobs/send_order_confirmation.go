package jobs

import (
	"context"
	"log/slog"

	"github.com/riverqueue/river"
)

// SendOrderConfirmationArgs is what a send_order_confirmation job carries. It is stored as JSON in the job row,
// so keep it to ids and small values; the worker loads what it needs.
type SendOrderConfirmationArgs struct {
	OrderID int64 `json:"order_id"`
}

// Kind names the job in the database. Never change it once jobs are queued;
// a queued job of the old kind would have no worker.
func (SendOrderConfirmationArgs) Kind() string { return "send_order_confirmation" }

// confirmer is what the worker needs: something that sends an order's email.
type confirmer interface {
	Send(ctx context.Context, orderID int64) error
}

// SendOrderConfirmationWorker does the work. Enqueue one with:
//
//	app.Jobs.Insert(ctx, jobs.SendOrderConfirmationArgs{}, nil)
//
// Dependencies go here as fields, declared as the smallest interface this
// worker uses, next to the worker, and are passed on the registration line in jobs.go.
type SendOrderConfirmationWorker struct {
	river.WorkerDefaults[SendOrderConfirmationArgs]
	Confirm confirmer
	Log     *slog.Logger
}

// Work runs one job. Return an error to have River retry it with backoff;
// return nil when it is done. ctx is cancelled when the worker drains.
func (w *SendOrderConfirmationWorker) Work(ctx context.Context, job *river.Job[SendOrderConfirmationArgs]) error {
	return w.Confirm.Send(ctx, job.Args.OrderID)
}
