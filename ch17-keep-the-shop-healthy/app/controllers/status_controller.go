package controllers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"depot/app/models"
)

type statusStore interface {
	Ping(ctx context.Context) error
	QueueHealth(ctx context.Context) (models.QueueHealth, error)
}

// StatusResponse is what /status reports.
type StatusResponse struct {
	Status         string `json:"status"`
	Database       string `json:"database"`
	WaitingTooLong int64  `json:"jobs_waiting_too_long"`
	FailedToday    int64  `json:"jobs_failed_today"`
}

// StatusController is the HTTP surface of status: whether the shop can do its
// work right now.
type StatusController struct {
	Store statusStore
	Log   *slog.Logger
}

// NewStatusController builds the controller.
func NewStatusController(store statusStore, log *slog.Logger) *StatusController {
	return &StatusController{Store: store, Log: log}
}

// SetupRoutes mounts the controller under its own prefix. Middleware the
// whole controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by the routes.go that calls this.
func (ctl *StatusController) SetupRoutes(router *gin.RouterGroup) {
	status := router.Group("/status")
	status.GET("", ctl.Index)
}

// Index handles GET /status: whether the shop can do its work right now.
func (ctl *StatusController) Index(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	resp := StatusResponse{Status: "ok", Database: "ok"}
	if err := ctl.Store.Ping(ctx); err != nil {
		ctl.Log.ErrorContext(ctx, "status: database unreachable", "err", err, "request_id", c.GetString("request_id"))
		resp.Status, resp.Database = "attention", "unreachable"
	} else if q, err := ctl.Store.QueueHealth(ctx); err != nil {
		ctl.Log.ErrorContext(ctx, "status: queue unreadable", "err", err, "request_id", c.GetString("request_id"))
		resp.Status = "attention"
	} else {
		resp.WaitingTooLong, resp.FailedToday = q.WaitingTooLong, q.FailedToday
		if q.WaitingTooLong > 0 || q.FailedToday > 0 {
			resp.Status = "attention"
		}
	}

	code := http.StatusOK
	if resp.Status != "ok" {
		code = http.StatusServiceUnavailable
	}
	c.JSON(code, resp)
}
