package jobs_test

import (
	"io"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/jobs"
)

var _ = Describe("Workers", func() {
	It("builds a bundle", func() {
		workers, _ := jobs.Workers(slog.New(slog.NewTextHandler(io.Discard, nil)), nil)
		Expect(workers).NotTo(BeNil())
	})
})
