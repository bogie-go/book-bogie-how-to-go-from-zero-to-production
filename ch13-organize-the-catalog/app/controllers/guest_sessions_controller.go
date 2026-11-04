package controllers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"depot/app/views"
)

// guestTokens is the slice of the token service this controller uses.
type guestTokens interface {
	Mint(subject string, ttl time.Duration) (string, error)
}

// guestSessionTTL is how long a guest's token, and so their cart, lasts.
const guestSessionTTL = 30 * 24 * time.Hour

// GuestSessionsController is the HTTP surface of guest_sessions: its dependencies, its routes
// and its actions, all in one file. Dependencies are small
// interfaces declared here, next to their consumer, holding only the methods the actions
// call, and the real thing is passed in from app/application.go, where this
// controller is wired:
//
//	type guestSessionsService interface {
//		Get(ctx context.Context, id string) (models.Thing, error)
//	}
type GuestSessionsController struct {
	Tokens guestTokens
	Log    *slog.Logger
}

// NewGuestSessionsController builds the controller.
func NewGuestSessionsController(tokens guestTokens, log *slog.Logger) *GuestSessionsController {
	return &GuestSessionsController{Tokens: tokens, Log: log}
}

// SetupRoutes mounts the controller under its own prefix. Middleware the
// whole controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by the routes.go that calls this.
func (ctl *GuestSessionsController) SetupRoutes(router *gin.RouterGroup) {
	guestSessions := router.Group("/guest_sessions")
	guestSessions.POST("", ctl.Create)
}

// Create handles POST /guest_sessions.
//
// A handler converts HTTP into models, calls one service function,
// and converts the result back. Nothing below app/controllers sees
// *gin.Context.
func (ctl *GuestSessionsController) Create(c *gin.Context) {
	signed, err := ctl.Tokens.Mint("guest:"+uuid.NewString(), guestSessionTTL)
	if err != nil {
		ctl.Log.ErrorContext(c.Request.Context(), "guest_sessions: mint failed", "err", err, "request_id", c.GetString("request_id"))
		c.JSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
		return
	}
	c.JSON(http.StatusCreated, views.GuestSessionResponse{Token: signed})
}
