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

// categoriesStore is the slice of the store this controller uses, declared next to
// its caller and holding only the methods the actions call.
// *adapter.Store satisfies it, and app/application.go passes it in.
type categoriesStore interface {
	ListCategorySummaries(ctx context.Context) ([]models.CategorySummary, error)
}

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
}

// Index handles GET /categories.
func (ctl *CategoriesController) Index(c *gin.Context) {
	categories, err := ctl.Store.ListCategorySummaries(c.Request.Context())
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, categoriesView(categories))
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
// imports models, so this is where CategoryResponse is built.

func categoryView(category models.CategorySummary) views.CategoryResponse {
	return views.CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		ProductCount: category.ProductCount,
	}
}

func categoriesView(categories []models.CategorySummary) views.CategoriesResponse {
	out := views.CategoriesResponse{Categories: make([]views.CategoryResponse, 0, len(categories))}
	for _, category := range categories {
		out.Categories = append(out.Categories, categoryView(category))
	}
	return out
}
