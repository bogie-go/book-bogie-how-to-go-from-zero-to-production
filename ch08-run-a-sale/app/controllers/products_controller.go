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

// productsStore is the slice of the store this controller uses, declared next to
// its caller and holding only the methods the actions call.
// *adapter.Store satisfies it, and app/application.go passes it in.
type productsStore interface {
	GetProduct(ctx context.Context, id string) (models.Product, error)
	ListProducts(ctx context.Context, limit int32) ([]models.Product, error)
}

// productsListLimit is how many a listing returns, newest first. Paging is
// yours to add when a client needs more.
const productsListLimit = 100

// ProductsController is the HTTP surface of products: a JSON resource over the store. A
// handler converts HTTP into models, calls one store method, and
// converts the result back. Validation is the model's and runs on save,
// which is why a bad value is a 422 from the store and not a binding rule
// here. Nothing below app/controllers sees *gin.Context.
type ProductsController struct {
	Store productsStore
	Log   *slog.Logger
}

// NewProductsController builds the controller.
func NewProductsController(store productsStore, log *slog.Logger) *ProductsController {
	return &ProductsController{Store: store, Log: log}
}

// SetupRoutes mounts the resource under its own prefix. Middleware the whole
// controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by routes.go, which calls this.
func (ctl *ProductsController) SetupRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	products.GET("", ctl.Index)
	products.GET("/:id", ctl.Show)
}

// Index handles GET /products.
func (ctl *ProductsController) Index(c *gin.Context) {
	products, err := ctl.Store.ListProducts(c.Request.Context(), productsListLimit)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, productsView(products))
}

// Show handles GET /products/:id.
func (ctl *ProductsController) Show(c *gin.Context) {
	p, err := ctl.Store.GetProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, productView(p))
}

// fail maps a model error to a status, the one place that mapping lives:
// not found is 404, a validation failure is 422, and anything
// else is logged and answered 500 without leaking its text.
func (ctl *ProductsController) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		c.JSON(http.StatusNotFound, views.ErrorMessage("no such product"))
	case errors.Is(err, models.ErrInvalid):
		c.JSON(http.StatusUnprocessableEntity, views.Error(err))
	default:
		ctl.Log.ErrorContext(c.Request.Context(), "products: request failed",
			"path", c.FullPath(), "err", err, "request_id", c.GetString("request_id"))
		c.JSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
	}
}

// The controller converts between the models and the wire. views never
// imports models, so this is where ProductResponse is built.

func productView(p models.Product) views.ProductResponse {
	return views.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		PriceCents:  p.PriceCents,
		Stock:       p.Stock,
		ImageUrl:    p.ImageUrl,
	}
}

func productsView(products []models.Product) views.ProductsResponse {
	out := views.ProductsResponse{Products: make([]views.ProductResponse, 0, len(products))}
	for _, p := range products {
		out.Products = append(out.Products, productView(p))
	}
	return out
}
