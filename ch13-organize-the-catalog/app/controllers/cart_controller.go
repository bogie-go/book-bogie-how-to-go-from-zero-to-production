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

// cartReader is what every cart response needs from the store.
type cartReader interface {
	CartContents(ctx context.Context, cartID int64) ([]models.CartLine, error)
	ListActivePriceRules(ctx context.Context) ([]models.PriceRule, error)
}

type cartStore interface {
	cartReader
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
func respondWithCart(c *gin.Context, store cartReader, log *slog.Logger, cartID int64) {
	ctx := c.Request.Context()
	lines, err := store.CartContents(ctx, cartID)
	if err != nil {
		failCart(c, log, err)
		return
	}
	rules, err := store.ListActivePriceRules(ctx)
	if err != nil {
		failCart(c, log, err)
		return
	}
	priced := models.PriceCart(lines, rules)
	for _, skipped := range priced.Skipped {
		log.WarnContext(ctx, "cart: skipped a price rule", "err", skipped, "request_id", c.GetString("request_id"))
	}
	c.JSON(http.StatusOK, cartView(lines, priced))
}

// cartView turns cart lines and their prices into the cart's JSON.
func cartView(lines []models.CartLine, priced models.PricedCart) views.CartResponse {
	resp := views.CartResponse{
		Items:         []views.CartLineResponse{},
		SubtotalCents: priced.SubtotalCents,
		DiscountCents: priced.DiscountCents(),
		TotalCents:    priced.TotalCents,
		AppliedRules:  []views.AppliedRuleResponse{},
	}
	for _, l := range lines {
		resp.Items = append(resp.Items, views.CartLineResponse{
			ID: l.ItemID, ProductID: l.ProductID, Name: l.Name, ImageUrl: l.ImageUrl,
			UnitPriceCents: l.PriceCents, Quantity: l.Quantity, SubtotalCents: l.PriceCents * l.Quantity,
		})
	}
	for _, a := range priced.Applied {
		resp.AppliedRules = append(resp.AppliedRules, views.AppliedRuleResponse{
			Name:             a.Name,
			PriceBeforeCents: a.BeforeCents,
			PriceAfterCents:  a.AfterCents,
			DiscountCents:    a.BeforeCents - a.AfterCents,
		})
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
