package controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"depot/app/models"
	"depot/app/views"
)

// ProductsController is the HTTP surface of products: its dependencies, its routes
// and its actions, all in one file. Dependencies are small
// interfaces declared here, next to their consumer, holding only the methods the actions
// call, and the real thing is passed in from app/application.go, where this
// controller is wired:
//
//	type productsService interface {
//		Get(ctx context.Context, id string) (models.Thing, error)
//	}
type ProductsController struct {
	Store productsStore
	Log   *slog.Logger
}

// productsStore is what this controller needs from the store.
type productsStore interface {
	ListProducts(ctx context.Context, limit int32) ([]models.Product, error)
	GetProduct(ctx context.Context, id string) (models.Product, error)
}

// productsListLimit is how many products a listing returns, newest first.
const productsListLimit = 100

// NewProductsController builds the controller.
func NewProductsController(store productsStore, log *slog.Logger) *ProductsController {
	return &ProductsController{Store: store, Log: log}
}

// SetupRoutes mounts the controller under its own prefix. Middleware the
// whole controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by the routes.go that calls this.
func (ctl *ProductsController) SetupRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	products.GET("", ctl.Index)
	products.GET("/:id", ctl.Show)
}

// Index handles GET /products.
//
// A handler converts HTTP into models, calls one service function,
// and converts the result back. Nothing below app/controllers sees
// *gin.Context.
func (ctl *ProductsController) Index(c *gin.Context) {
	products, err := ctl.Store.ListProducts(c.Request.Context(), productsListLimit)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	resp := views.ProductsResponse{Products: make([]views.ProductResponse, 0, len(products))}
	for _, p := range products {
		resp.Products = append(resp.Products, productView(p))
	}
	c.JSON(http.StatusOK, resp)
}

// Show handles GET /products/:id.
//
// A handler converts HTTP into models, calls one service function,
// and converts the result back. Nothing below app/controllers sees
// *gin.Context.
func (ctl *ProductsController) Show(c *gin.Context) {
	p, err := ctl.Store.GetProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, productView(p))
}

// productView turns a product into what the shop shows of it.
func productView(p models.Product) views.ProductResponse {
	return views.ProductResponse{
		ID: p.ID, Name: p.Name, Description: p.Description,
		PriceCents: p.PriceCents, Stock: p.Stock, ImageUrl: p.ImageUrl,
	}
}

// fail is the one place an error becomes a status.
func (ctl *ProductsController) fail(c *gin.Context, err error) {
	if errors.Is(err, models.ErrNotFound) {
		c.JSON(http.StatusNotFound, views.ErrorMessage("no such product"))
		return
	}
	ctl.Log.ErrorContext(c.Request.Context(), "products: request failed",
		"path", c.FullPath(), "err", err, "request_id", c.GetString("request_id"))
	c.JSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
}
