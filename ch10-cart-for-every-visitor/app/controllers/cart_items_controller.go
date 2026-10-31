package controllers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"depot/app/middlewares"
	"depot/app/models"
	"depot/app/views"
)

type cartItemsStore interface {
	cartContents
	FindOrCreateCartByToken(ctx context.Context, token string) (models.Cart, error)
	AddCartItem(ctx context.Context, cartID, productID int64, quantity int32) error
	SetCartItemQuantity(ctx context.Context, cartID int64, id string, quantity int32) (models.CartItem, error)
	RemoveCartItem(ctx context.Context, cartID int64, id string) error
}

// CartItemsController is the HTTP surface of cart_items: its dependencies, its routes
// and its actions, all in one file. Dependencies are small
// interfaces declared here, next to their consumer, holding only the methods the actions
// call, and the real thing is passed in from app/application.go, where this
// controller is wired:
//
//	type cartItemsService interface {
//		Get(ctx context.Context, id string) (models.Thing, error)
//	}
type CartItemsController struct {
	Store cartItemsStore
	Log   *slog.Logger
}

// NewCartItemsController builds the controller.
func NewCartItemsController(store cartItemsStore, log *slog.Logger) *CartItemsController {
	return &CartItemsController{Store: store, Log: log}
}

// SetupRoutes mounts the controller under its own prefix. Middleware the
// whole controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by the routes.go that calls this.
func (ctl *CartItemsController) SetupRoutes(router *gin.RouterGroup) {
	cartItems := router.Group("/cart_items")
	cartItems.POST("", ctl.Create)
	cartItems.PATCH("/:id", ctl.Update)
	cartItems.DELETE("/:id", ctl.Destroy)
}

// Create handles POST /cart_items.
//
// A handler converts HTTP into models, calls one service function,
// and converts the result back. Nothing below app/controllers sees
// *gin.Context.
func (ctl *CartItemsController) Create(c *gin.Context) {
	var req views.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorMessage("invalid JSON"))
		return
	}
	ctx := c.Request.Context()
	cart, err := ctl.Store.FindOrCreateCartByToken(ctx, middlewares.Subject(c))
	if err != nil {
		failCart(c, ctl.Log, err)
		return
	}
	if err := ctl.Store.AddCartItem(ctx, cart.ID, req.ProductID, req.Quantity); err != nil {
		failCart(c, ctl.Log, err)
		return
	}
	respondWithCart(c, ctl.Store, ctl.Log, cart.ID)
}

// Update handles PATCH /cart_items/:id.
//
// A handler converts HTTP into models, calls one service function,
// and converts the result back. Nothing below app/controllers sees
// *gin.Context.
func (ctl *CartItemsController) Update(c *gin.Context) {
	var req views.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorMessage("invalid JSON"))
		return
	}
	ctx := c.Request.Context()
	cart, err := ctl.Store.FindOrCreateCartByToken(ctx, middlewares.Subject(c))
	if err != nil {
		failCart(c, ctl.Log, err)
		return
	}
	if _, err := ctl.Store.SetCartItemQuantity(ctx, cart.ID, c.Param("id"), req.Quantity); err != nil {
		failCart(c, ctl.Log, err)
		return
	}
	respondWithCart(c, ctl.Store, ctl.Log, cart.ID)
}

// Destroy handles DELETE /cart_items/:id.
//
// A handler converts HTTP into models, calls one service function,
// and converts the result back. Nothing below app/controllers sees
// *gin.Context.
func (ctl *CartItemsController) Destroy(c *gin.Context) {
	ctx := c.Request.Context()
	cart, err := ctl.Store.FindOrCreateCartByToken(ctx, middlewares.Subject(c))
	if err != nil {
		failCart(c, ctl.Log, err)
		return
	}
	if err := ctl.Store.RemoveCartItem(ctx, cart.ID, c.Param("id")); err != nil {
		failCart(c, ctl.Log, err)
		return
	}
	respondWithCart(c, ctl.Store, ctl.Log, cart.ID)
}
