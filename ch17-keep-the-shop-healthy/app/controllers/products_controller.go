package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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
	ListCatalogProducts(ctx context.Context, categoryID *int64, limit, offset int32) ([]models.Product, error)
	CountCatalogProducts(ctx context.Context, categoryID *int64) (int64, error)
	CategoriesForProduct(ctx context.Context, productID int64) ([]models.Category, error)
	GetCatalogProduct(ctx context.Context, id string) (models.Product, error)
}

const (
	defaultPerPage = 12
	maxPerPage     = 48
)

// pageParam reads a positive whole number from the query, or def when it's
// missing.
func pageParam(c *gin.Context, name string, def int) (int, bool) {
	raw := c.Query(name)
	if raw == "" {
		return def, true
	}
	n, err := strconv.Atoi(raw)
	return n, err == nil && n >= 1
}

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

// Index handles GET /products, optionally ?category_id=N&page=N&per_page=N.
//
// A handler converts HTTP into models, calls one service function,
// and converts the result back. Nothing below app/controllers sees
// *gin.Context.
func (ctl *ProductsController) Index(c *gin.Context) {
	var categoryID *int64
	if raw := c.Query("category_id"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, views.ErrorMessage("category_id must be a number"))
			return
		}
		categoryID = &id
	}
	page, okPage := pageParam(c, "page", 1)
	perPage, okPer := pageParam(c, "per_page", defaultPerPage)
	if !okPage || !okPer || perPage > maxPerPage {
		c.JSON(http.StatusBadRequest, views.ErrorMessage(fmt.Sprintf("page and per_page must be whole numbers from 1, and per_page at most %d", maxPerPage)))
		return
	}
	ctx := c.Request.Context()
	products, err := ctl.Store.ListCatalogProducts(ctx, categoryID, int32(perPage), int32((page-1)*perPage))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	total, err := ctl.Store.CountCatalogProducts(ctx, categoryID)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	resp := views.ProductsResponse{Products: make([]views.ProductResponse, 0, len(products)), Page: page, PerPage: perPage, Total: total}
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
	ctx := c.Request.Context()
	p, err := ctl.Store.GetCatalogProduct(ctx, c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	categories, err := ctl.Store.CategoriesForProduct(ctx, p.ID)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	resp := views.ProductDetailResponse{
		ProductResponse: productView(p),
		Categories:      make([]views.CategoryNameResponse, 0, len(categories)),
	}
	for _, cat := range categories {
		resp.Categories = append(resp.Categories, views.CategoryNameResponse{ID: cat.ID, Name: cat.Name})
	}
	c.JSON(http.StatusOK, resp)
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
