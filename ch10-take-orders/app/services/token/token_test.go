package token_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/services/token"
)

var maker = token.Maker{Secret: []byte("a-test-secret-of-at-least-thirty-two-bytes"), Issuer: "test"}

var _ = Describe("Maker", func() {
	var (
		raw    string
		claims token.Claims
		err    error
	)

	mint := func(m token.Maker, ttl time.Duration) string {
		signed, err := m.Mint("user-42", ttl)
		Expect(err).NotTo(HaveOccurred())
		return signed
	}

	Describe("Mint and Verify", func() {
		Context("with a valid token", func() {
			BeforeEach(func() {
				raw = mint(maker, time.Minute)
				claims, err = maker.Verify(raw)
			})

			It("verifies", func() { Expect(err).NotTo(HaveOccurred()) })
			It("gives back the subject", func() { Expect(claims.Subject).To(Equal("user-42")) })
			It("carries no role unless one was minted", func() { Expect(claims.Role).To(BeEmpty()) })
			It("carries an id", func() { Expect(claims.ID).NotTo(BeEmpty()) })
			It("expires in the future", func() { Expect(claims.ExpiresAt).To(BeTemporally(">", time.Now())) })
		})

		Context("with an expired token", func() {
			BeforeEach(func() {
				_, err = maker.Verify(mint(maker, -time.Minute))
			})

			It("is refused as expired", func() { Expect(err).To(MatchError(token.ErrExpired)) })
		})

		Context("with a token signed by another secret", func() {
			BeforeEach(func() {
				other := token.Maker{Secret: []byte("another-secret-of-at-least-thirty-two-bytes"), Issuer: "test"}
				_, err = maker.Verify(mint(other, time.Minute))
			})

			It("is refused as invalid", func() { Expect(err).To(MatchError(token.ErrInvalid)) })
		})

		Context("with a token from another issuer", func() {
			BeforeEach(func() {
				other := token.Maker{Secret: maker.Secret, Issuer: "someone-else"}
				_, err = maker.Verify(mint(other, time.Minute))
			})

			It("is refused as invalid", func() { Expect(err).To(MatchError(token.ErrInvalid)) })
		})

		Context("with garbage", func() {
			BeforeEach(func() {
				_, err = maker.Verify("not.a.token")
			})

			It("is refused as invalid", func() { Expect(err).To(MatchError(token.ErrInvalid)) })
		})

		Context("with an empty string", func() {
			BeforeEach(func() {
				_, err = maker.Verify("")
			})

			It("is refused as invalid", func() { Expect(err).To(MatchError(token.ErrInvalid)) })
		})
	})

	Describe("Mint", func() {
		Context("with an empty subject", func() {
			BeforeEach(func() {
				_, err = maker.Mint("", time.Minute)
			})

			It("refuses to mint", func() { Expect(err).To(HaveOccurred()) })
		})
	})

	Describe("MintFor", func() {
		Context("with a role", func() {
			BeforeEach(func() {
				raw, err = maker.MintFor("user-42", "admin", time.Minute)
				Expect(err).NotTo(HaveOccurred())
				claims, err = maker.Verify(raw)
			})

			It("keeps the role", func() { Expect(claims.Role).To(Equal("admin")) })
			It("keeps the subject", func() { Expect(claims.Subject).To(Equal("user-42")) })
		})
	})
})
