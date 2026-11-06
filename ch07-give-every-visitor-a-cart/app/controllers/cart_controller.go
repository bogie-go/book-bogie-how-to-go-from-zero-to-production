package controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"depot/app/middlewares"
	"depot/app/models"
	"depot/app/views"
)

// cartContents is the one store method every cart response needs.
type cartContents interface {
	CartContents(ctx context.Context, cartID int64) ([]models.CartLine, error)
}

type cartStore interface {
	cartContents
	FindOrCreateCartByToken(ctx context.Context, token string) (models.Cart, error)
}

// CartController is the HTTP surface of cart: its dependencies, its routes
// and its actions, all in one file. Dependencies are small
// interfaces declared here, next to their consumer, holding only the methods the actions
// call, and the real thing is passed in from app/application.go, where this
// controller is wired:
//
//	type cartService interface {
//		Get(ctx context.Context, id string) (models.Thing, error)
//	}
type CartController struct {
	Store cartStore
	Log   *slog.Logger
}

// NewCartController builds the controller.
func NewCartController(store cartStore, log *slog.Logger) *CartController {
	return &CartController{Store: store, Log: log}
}

// SetupRoutes mounts the controller under its own prefix. Middleware the
// whole controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by the routes.go that calls this.
func (ctl *CartController) SetupRoutes(router *gin.RouterGroup) {
	cart := router.Group("/cart")
	cart.GET("", ctl.Index)
}

// Index handles GET /cart.
//
// A handler converts HTTP into models, calls one service function,
// and converts the result back. Nothing below app/controllers sees
// *gin.Context.
func (ctl *CartController) Index(c *gin.Context) {
	cart, err := ctl.Store.FindOrCreateCartByToken(c.Request.Context(), middlewares.Subject(c))
	if err != nil {
		failCart(c, ctl.Log, err)
		return
	}
	respondWithCart(c, ctl.Store, ctl.Log, cart.ID)
}

// respondWithCart answers with the whole cart: every change answers this way.
func respondWithCart(c *gin.Context, store cartContents, log *slog.Logger, cartID int64) {
	lines, err := store.CartContents(c.Request.Context(), cartID)
	if err != nil {
		failCart(c, log, err)
		return
	}
	c.JSON(http.StatusOK, cartView(lines))
}

// cartView turns cart lines into the cart's JSON, working out the subtotals
// and the total as it goes.
func cartView(lines []models.CartLine) views.CartResponse {
	resp := views.CartResponse{Items: []views.CartLineResponse{}}
	for _, l := range lines {
		subtotal := l.PriceCents * l.Quantity
		resp.Items = append(resp.Items, views.CartLineResponse{
			ID: l.ItemID, ProductID: l.ProductID, Name: l.Name, ImageUrl: l.ImageUrl,
			UnitPriceCents: l.PriceCents, Quantity: l.Quantity, SubtotalCents: subtotal,
		})
		resp.TotalCents += subtotal
	}
	return resp
}

// failCart turns a cart error into the right status.
func failCart(c *gin.Context, log *slog.Logger, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		c.JSON(http.StatusNotFound, views.ErrorMessage("not found"))
	case errors.Is(err, models.ErrInvalid):
		c.JSON(http.StatusUnprocessableEntity, views.Error(err))
	default:
		log.ErrorContext(c.Request.Context(), "cart: request failed", "path", c.FullPath(), "err", err, "request_id", c.GetString("request_id"))
		c.JSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
	}
}
