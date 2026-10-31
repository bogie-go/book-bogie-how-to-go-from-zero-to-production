package controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"depot/app/models"
	"depot/app/services/password"
	"depot/app/views"
)

// userSessionTTL is how long a customer stays signed in.
const userSessionTTL = 7 * 24 * time.Hour

type usersStore interface {
	cartMerger
	CreateUser(ctx context.Context, u models.User) (models.User, error)
}

// UsersController is the HTTP surface of users: its dependencies, its routes
// and its actions, all in one file.
type UsersController struct {
	Store  usersStore
	Tokens accountTokens
	Log    *slog.Logger
}

// NewUsersController builds the controller.
func NewUsersController(store usersStore, tokens accountTokens, log *slog.Logger) *UsersController {
	return &UsersController{Store: store, Tokens: tokens, Log: log}
}

// SetupRoutes mounts the controller under its own prefix. Middleware the
// whole controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by the routes.go that calls this.
func (ctl *UsersController) SetupRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.POST("", ctl.Create)
}

// Create handles POST /users: a new account, signed in straight away.
func (ctl *UsersController) Create(c *gin.Context) {
	var req views.CredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorMessage("invalid JSON"))
		return
	}
	digest, err := password.Hash(req.Password)
	if err != nil {
		failAccount(c, ctl.Log, err)
		return
	}
	user, err := ctl.Store.CreateUser(c.Request.Context(), models.User{
		Email:          models.NormalizeEmail(req.Email),
		PasswordDigest: digest,
	})
	if err != nil {
		failAccount(c, ctl.Log, err)
		return
	}
	subject := models.UserSubject(user.ID)
	mergeGuestCart(c, ctl.Tokens, ctl.Store, ctl.Log, subject)
	signed, err := ctl.Tokens.Mint(subject, userSessionTTL)
	if err != nil {
		failAccount(c, ctl.Log, err)
		return
	}
	c.JSON(http.StatusCreated, views.SessionResponse{Token: signed, Email: user.Email})
}

// failAccount turns an account error into the right status.
func failAccount(c *gin.Context, log *slog.Logger, err error) {
	if errors.Is(err, models.ErrInvalid) {
		c.JSON(http.StatusUnprocessableEntity, views.Error(err))
		return
	}
	log.ErrorContext(c.Request.Context(), "accounts: request failed", "path", c.FullPath(), "err", err, "request_id", c.GetString("request_id"))
	c.JSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
}
