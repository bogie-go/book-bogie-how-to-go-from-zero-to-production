package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/middlewares"
	"depot/app/services/token"
)

var testTokens = token.Maker{Secret: []byte("a-test-secret-of-at-least-thirty-two-bytes"), Issuer: "test"}

func userRequest(authorization string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/me", middlewares.RequireUser(testTokens), func(c *gin.Context) { c.String(http.StatusOK, middlewares.Subject(c)) })
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

var _ = Describe("RequireUser", func() {
	var rec *httptest.ResponseRecorder

	Context("with a good token", func() {
		BeforeEach(func() {
			good, err := testTokens.Mint("user-42", time.Minute)
			Expect(err).NotTo(HaveOccurred())
			rec = userRequest("Bearer " + good)
		})

		It("lets the request through", func() { Expect(rec.Code).To(Equal(http.StatusOK)) })
		It("puts the subject on the context", func() { Expect(rec.Body.String()).To(Equal("user-42")) })
	})

	Context("with no token", func() {
		BeforeEach(func() { rec = userRequest("") })

		It("answers 401", func() { Expect(rec.Code).To(Equal(http.StatusUnauthorized)) })
	})

	Context("with an expired token", func() {
		BeforeEach(func() {
			expired, err := testTokens.Mint("user-42", -time.Minute)
			Expect(err).NotTo(HaveOccurred())
			rec = userRequest("Bearer " + expired)
		})

		It("answers 401", func() { Expect(rec.Code).To(Equal(http.StatusUnauthorized)) })
	})

	Context("with garbage", func() {
		BeforeEach(func() { rec = userRequest("Bearer nope") })

		It("answers 401", func() { Expect(rec.Code).To(Equal(http.StatusUnauthorized)) })
	})

	Context("with the wrong scheme", func() {
		BeforeEach(func() {
			good, err := testTokens.Mint("user-42", time.Minute)
			Expect(err).NotTo(HaveOccurred())
			rec = userRequest("Basic " + good)
		})

		It("answers 401", func() { Expect(rec.Code).To(Equal(http.StatusUnauthorized)) })
	})
})
