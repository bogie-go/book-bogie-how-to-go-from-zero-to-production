package controllers_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("POST /guest_sessions", func() {
	It("answers 201 with a new token", func() {
		app.Post("/guest_sessions", nil).Do().
			Status(http.StatusCreated).
			Contains(`"token":"`)
	})
})
