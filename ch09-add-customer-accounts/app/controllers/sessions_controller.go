package controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"depot/app/models"
	"depot/app/services/password"
	"depot/app/views"
)

type sessionsStore interface {
	cartMerger
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
}

// SessionsController is the HTTP surface of sessions: its dependencies, its routes
// and its actions, all in one file.
type SessionsController struct {
	Store  sessionsStore
	Tokens accountTokens
	Log    *slog.Logger
}

// NewSessionsController builds the controller.
func NewSessionsController(store sessionsStore, tokens accountTokens, log *slog.Logger) *SessionsController {
	return &SessionsController{Store: store, Tokens: tokens, Log: log}
}

// SetupRoutes mounts the controller under its own prefix. Middleware the
// whole controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by the routes.go that calls this.
func (ctl *SessionsController) SetupRoutes(router *gin.RouterGroup) {
	sessions := router.Group("/sessions")
	sessions.POST("", ctl.Create)
}

// Create handles POST /sessions: sign in with an email and a password.
func (ctl *SessionsController) Create(c *gin.Context) {
	var req views.CredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorMessage("invalid JSON"))
		return
	}
	user, err := ctl.Store.FindUserByEmail(c.Request.Context(), models.NormalizeEmail(req.Email))
	switch {
	case errors.Is(err, models.ErrNotFound):
		password.DummyCheck(req.Password)
		refuseSignIn(c)
		return
	case err != nil:
		failAccount(c, ctl.Log, err)
		return
	}
	if !password.Check(user.PasswordDigest, req.Password) {
		refuseSignIn(c)
		return
	}
	subject := userSubject(user.ID)
	mergeGuestCart(c, ctl.Tokens, ctl.Store, ctl.Log, subject)
	signed, err := ctl.Tokens.Mint(subject, userSessionTTL)
	if err != nil {
		failAccount(c, ctl.Log, err)
		return
	}
	c.JSON(http.StatusOK, views.SessionResponse{Token: signed, Email: user.Email})
}

// refuseSignIn is the one answer to every wrong email or password.
func refuseSignIn(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, views.ErrorMessage("email or password is incorrect"))
}
