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

// priceRulesStore is the slice of the store this controller uses, declared next to
// its caller and holding only the methods the actions call.
// *adapter.Store satisfies it, and app/application.go passes it in.
type priceRulesStore interface {
	CreatePriceRule(ctx context.Context, p models.PriceRule) (models.PriceRule, error)
	GetPriceRule(ctx context.Context, id string) (models.PriceRule, error)
	ListPriceRules(ctx context.Context, limit int32) ([]models.PriceRule, error)
	UpdatePriceRule(ctx context.Context, p models.PriceRule) (models.PriceRule, error)
	DeletePriceRule(ctx context.Context, id string) error
}

// priceRulesListLimit is how many a listing returns, newest first. Paging is
// yours to add when a client needs more.
const priceRulesListLimit = 100

// PriceRulesController is the HTTP surface of price_rules: a JSON resource over the store. A
// handler converts HTTP into models, calls one store method, and
// converts the result back. Validation is the model's and runs on save,
// which is why a bad value is a 422 from the store and not a binding rule
// here. Nothing below app/controllers sees *gin.Context.
type PriceRulesController struct {
	Store priceRulesStore
	Log   *slog.Logger
}

// NewPriceRulesController builds the controller.
func NewPriceRulesController(store priceRulesStore, log *slog.Logger) *PriceRulesController {
	return &PriceRulesController{Store: store, Log: log}
}

// SetupRoutes mounts the resource under its own prefix. Middleware the whole
// controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by routes.go, which calls this.
func (ctl *PriceRulesController) SetupRoutes(router *gin.RouterGroup) {
	priceRules := router.Group("/price_rules")
	priceRules.GET("", ctl.Index)
	priceRules.POST("", ctl.Create)
	priceRules.GET("/:id", ctl.Show)
	priceRules.PATCH("/:id", ctl.Update)
	priceRules.DELETE("/:id", ctl.Destroy)
}

// Index handles GET /admin/price_rules.
func (ctl *PriceRulesController) Index(c *gin.Context) {
	priceRules, err := ctl.Store.ListPriceRules(c.Request.Context(), priceRulesListLimit)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, priceRulesView(priceRules))
}

// Show handles GET /admin/price_rules/:id.
func (ctl *PriceRulesController) Show(c *gin.Context) {
	p, err := ctl.Store.GetPriceRule(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, priceRuleView(p))
}

// Create handles POST /admin/price_rules.
func (ctl *PriceRulesController) Create(c *gin.Context) {
	var req views.AdminPriceRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.Error(err))
		return
	}
	p, err := ctl.Store.CreatePriceRule(c.Request.Context(), models.PriceRule{Name: req.Name, Kind: req.Kind, Properties: req.Properties, Priority: req.Priority, StartsAt: req.StartsAt, EndsAt: req.EndsAt})
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, priceRuleView(p))
}

// Update handles PATCH /admin/price_rules/:id. Only the fields in the body change: the price_rule is loaded,
// the body is decoded over its current values,
// so a key left out keeps its value, and the result is saved and validated
// whole.
func (ctl *PriceRulesController) Update(c *gin.Context) {
	p, err := ctl.Store.GetPriceRule(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	req := views.AdminPriceRuleRequest{Name: p.Name, Kind: p.Kind, Properties: p.Properties, Priority: p.Priority, StartsAt: p.StartsAt, EndsAt: p.EndsAt}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.Error(err))
		return
	}
	p, err = ctl.Store.UpdatePriceRule(c.Request.Context(), models.PriceRule{ID: p.ID, Name: req.Name, Kind: req.Kind, Properties: req.Properties, Priority: req.Priority, StartsAt: req.StartsAt, EndsAt: req.EndsAt})
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, priceRuleView(p))
}

// Destroy handles DELETE /admin/price_rules/:id.
func (ctl *PriceRulesController) Destroy(c *gin.Context) {
	if err := ctl.Store.DeletePriceRule(c.Request.Context(), c.Param("id")); err != nil {
		ctl.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// fail maps a model error to a status, the one place that mapping lives:
// not found is 404, a validation failure is 422, and anything
// else is logged and answered 500 without leaking its text.
func (ctl *PriceRulesController) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		c.JSON(http.StatusNotFound, views.ErrorMessage("no such price_rule"))
	case errors.Is(err, models.ErrInvalid):
		c.JSON(http.StatusUnprocessableEntity, views.Error(err))
	default:
		ctl.Log.ErrorContext(c.Request.Context(), "price_rules: request failed",
			"path", c.FullPath(), "err", err, "request_id", c.GetString("request_id"))
		c.JSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
	}
}

// The controller converts between the models and the wire. views never
// imports models, so this is where AdminPriceRuleResponse is built.

func priceRuleView(p models.PriceRule) views.AdminPriceRuleResponse {
	return views.AdminPriceRuleResponse{
		ID:         p.ID,
		Name:       p.Name,
		Kind:       p.Kind,
		Properties: p.Properties,
		Priority:   p.Priority,
		StartsAt:   p.StartsAt,
		EndsAt:     p.EndsAt,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func priceRulesView(priceRules []models.PriceRule) views.AdminPriceRulesResponse {
	out := views.AdminPriceRulesResponse{PriceRules: make([]views.AdminPriceRuleResponse, 0, len(priceRules))}
	for _, p := range priceRules {
		out.PriceRules = append(out.PriceRules, priceRuleView(p))
	}
	return out
}
