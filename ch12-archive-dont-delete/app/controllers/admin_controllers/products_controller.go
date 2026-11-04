package admin_controllers

import (
	"context"
	"encoding/json"
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
	CreateProduct(ctx context.Context, p models.Product) (models.Product, error)
	GetProduct(ctx context.Context, id string) (models.Product, error)
	ListProducts(ctx context.Context, limit int32) ([]models.Product, error)
	UpdateProduct(ctx context.Context, p models.Product) (models.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ArchiveProduct(ctx context.Context, id string) (models.Product, error)
	RestoreProduct(ctx context.Context, id string) (models.Product, error)
	EditProduct(ctx context.Context, id string, edit func(*models.Product) error) (models.Product, error)
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
	products.POST("", ctl.Create)
	products.GET("/:id", ctl.Show)
	products.PATCH("/:id", ctl.Update)
	products.DELETE("/:id", ctl.Destroy)
	products.POST("/:id/archive", ctl.Archive)
	products.POST("/:id/restore", ctl.Restore)
}

// Index handles GET /admin/products.
func (ctl *ProductsController) Index(c *gin.Context) {
	products, err := ctl.Store.ListProducts(c.Request.Context(), productsListLimit)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, productsView(products))
}

// Show handles GET /admin/products/:id.
func (ctl *ProductsController) Show(c *gin.Context) {
	p, err := ctl.Store.GetProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, productView(p))
}

// Create handles POST /admin/products.
func (ctl *ProductsController) Create(c *gin.Context) {
	var req views.AdminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.Error(err))
		return
	}
	p, err := ctl.Store.CreateProduct(c.Request.Context(), models.Product{Name: req.Name, Description: req.Description, PriceCents: req.PriceCents, Stock: req.Stock, ImageUrl: req.ImageUrl})
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, productView(p))
}

// Update handles PATCH /admin/products/:id. Only the fields in the body change:
// the body is decoded over the product's current values, read inside a lock
// on its row, so a key left out keeps its value and a checkout can't be undone.
func (ctl *ProductsController) Update(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil || !json.Valid(body) {
		c.JSON(http.StatusBadRequest, views.ErrorMessage("invalid JSON"))
		return
	}
	product, err := ctl.Store.EditProduct(c.Request.Context(), c.Param("id"), func(p *models.Product) error {
		req := views.AdminProductRequest{
			Name: p.Name, Description: p.Description, PriceCents: p.PriceCents,
			Stock: p.Stock, ImageUrl: p.ImageUrl,
		}
		if err := json.Unmarshal(body, &req); err != nil {
			return errBadBody
		}
		p.Name, p.Description, p.PriceCents = req.Name, req.Description, req.PriceCents
		p.Stock, p.ImageUrl = req.Stock, req.ImageUrl
		return nil
	})
	if errors.Is(err, errBadBody) {
		c.JSON(http.StatusBadRequest, views.ErrorMessage("invalid JSON"))
		return
	}
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, productView(product))
}

// errBadBody is a body that is JSON, but not a product's: a price as text, say.
var errBadBody = errors.New("invalid JSON")

// Destroy handles DELETE /admin/products/:id.
func (ctl *ProductsController) Destroy(c *gin.Context) {
	if err := ctl.Store.DeleteProduct(c.Request.Context(), c.Param("id")); err != nil {
		ctl.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Archive handles POST /admin/products/:id/archive: take a product out of the shop.
func (ctl *ProductsController) Archive(c *gin.Context) {
	product, err := ctl.Store.ArchiveProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, productView(product))
}

// Restore handles POST /admin/products/:id/restore: bring a product back.
func (ctl *ProductsController) Restore(c *gin.Context) {
	product, err := ctl.Store.RestoreProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, productView(product))
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
// imports models, so this is where AdminProductResponse is built.

func productView(p models.Product) views.AdminProductResponse {
	return views.AdminProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		PriceCents:  p.PriceCents,
		Stock:       p.Stock,
		ImageUrl:    p.ImageUrl,
		ArchivedAt:  p.ArchivedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func productsView(products []models.Product) views.AdminProductsResponse {
	out := views.AdminProductsResponse{Products: make([]views.AdminProductResponse, 0, len(products))}
	for _, p := range products {
		out.Products = append(out.Products, productView(p))
	}
	return out
}
