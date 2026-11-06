package controllers

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"depot/app/services/token"
)

// accountTokens mints a customer's token, and reads the one a visitor who's
// signing up or signing in already had.
type accountTokens interface {
	Mint(subject string, ttl time.Duration) (string, error)
	Verify(raw string) (token.Claims, error)
}

// cartMerger moves one cart's lines into another.
type cartMerger interface {
	MergeCart(ctx context.Context, from, to string) ([]error, error)
}

// mergeGuestCart brings the cart of the guest token this request carries, if
// any, into the customer's cart. Nothing about it can stop a sign-up or a
// sign-in.
func mergeGuestCart(c *gin.Context, tokens accountTokens, store cartMerger, log *slog.Logger, subject string) {
	raw, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
	if !ok {
		return
	}
	claims, err := tokens.Verify(raw)
	if err != nil || !strings.HasPrefix(claims.Subject, "guest:") {
		return
	}
	ctx := c.Request.Context()
	skipped, err := store.MergeCart(ctx, claims.Subject, subject)
	for _, s := range skipped {
		log.WarnContext(ctx, "accounts: a guest cart line didn't merge", "err", s, "request_id", c.GetString("request_id"))
	}
	if err != nil {
		log.WarnContext(ctx, "accounts: guest cart merge failed", "err", err, "request_id", c.GetString("request_id"))
	}
}
