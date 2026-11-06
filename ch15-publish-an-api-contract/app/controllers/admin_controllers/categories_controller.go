package admin_controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"depot/app/models"
	"depot/app/views"
)

// categoriesStore is the slice of the store this controller uses, declared next to
// its caller and holding only the methods the actions call.
// *adapter.Store satisfies it, and app/application.go passes it in.
type categoriesStore interface {
	CreateCategory(ctx context.Context, category models.Category) (models.Category, error)
	GetCategory(ctx context.Context, id string) (models.Category, error)
	ListCategories(ctx context.Context, limit int32) ([]models.Category, error)
	UpdateCategory(ctx context.Context, category models.Category) (models.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

// categoriesListLimit is how many a listing returns, newest first. Paging is
// yours to add when a client needs more.
const categoriesListLimit = 100

// CategoriesController is the HTTP surface of categories: a JSON resource over the store. A
// handler converts HTTP into models, calls one store method, and
// converts the result back. Validation is the model's and runs on save,
// which is why a bad value is a 422 from the store and not a binding rule
// here. Nothing below app/controllers sees *gin.Context.
type CategoriesController struct {
	Store categoriesStore
	Log   *slog.Logger
}

// NewCategoriesController builds the controller.
func NewCategoriesController(store categoriesStore, log *slog.Logger) *CategoriesController {
	return &CategoriesController{Store: store, Log: log}
}

// SetupRoutes mounts the resource under its own prefix. Middleware the whole
// controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by routes.go, which calls this.
func (ctl *CategoriesController) SetupRoutes(router *gin.RouterGroup) {
	categories := router.Group("/categories")
	categories.GET("", ctl.Index)
	categories.POST("", ctl.Create)
	categories.GET("/:id", ctl.Show)
	categories.PATCH("/:id", ctl.Update)
	categories.DELETE("/:id", ctl.Destroy)
}

// Index handles GET /admin/categories.
func (ctl *CategoriesController) Index(c *gin.Context) {
	categories, err := ctl.Store.ListCategories(c.Request.Context(), categoriesListLimit)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, categoriesView(categories))
}

// Show handles GET /admin/categories/:id.
func (ctl *CategoriesController) Show(c *gin.Context) {
	category, err := ctl.Store.GetCategory(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, categoryView(category))
}

// Create handles POST /admin/categories.
func (ctl *CategoriesController) Create(c *gin.Context) {
	var req views.AdminCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.Error(err))
		return
	}
	category, err := ctl.Store.CreateCategory(c.Request.Context(), models.Category{Name: req.Name})
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, categoryView(category))
}

// Update handles PATCH /admin/categories/:id. Only the fields in the body change: the category is loaded,
// the body is decoded over its current values,
// so a key left out keeps its value, and the result is saved and validated
// whole.
func (ctl *CategoriesController) Update(c *gin.Context) {
	category, err := ctl.Store.GetCategory(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	req := views.AdminCategoryRequest{Name: category.Name}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.Error(err))
		return
	}
	category, err = ctl.Store.UpdateCategory(c.Request.Context(), models.Category{ID: category.ID, Name: req.Name})
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, categoryView(category))
}

// Destroy handles DELETE /admin/categories/:id.
func (ctl *CategoriesController) Destroy(c *gin.Context) {
	if err := ctl.Store.DeleteCategory(c.Request.Context(), c.Param("id")); err != nil {
		ctl.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// fail maps a model error to a status, the one place that mapping lives:
// not found is 404, a validation failure is 422, and anything
// else is logged and answered 500 without leaking its text.
func (ctl *CategoriesController) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		c.JSON(http.StatusNotFound, views.ErrorMessage("no such category"))
	case errors.Is(err, models.ErrInvalid):
		c.JSON(http.StatusUnprocessableEntity, views.Error(err))
	default:
		ctl.Log.ErrorContext(c.Request.Context(), "categories: request failed",
			"path", c.FullPath(), "err", err, "request_id", c.GetString("request_id"))
		c.JSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
	}
}

// The controller converts between the models and the wire. views never
// imports models, so this is where AdminCategoryResponse is built.

func categoryView(category models.Category) views.AdminCategoryResponse {
	return views.AdminCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func categoriesView(categories []models.Category) views.AdminCategoriesResponse {
	out := views.AdminCategoriesResponse{Categories: make([]views.AdminCategoryResponse, 0, len(categories))}
	for _, category := range categories {
		out.Categories = append(out.Categories, categoryView(category))
	}
	return out
}
