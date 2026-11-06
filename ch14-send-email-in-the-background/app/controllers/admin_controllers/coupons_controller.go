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

// couponsStore is the slice of the store this controller uses, declared next to
// its caller and holding only the methods the actions call.
// *adapter.Store satisfies it, and app/application.go passes it in.
type couponsStore interface {
	CreateCoupon(ctx context.Context, coupon models.Coupon) (models.Coupon, error)
	GetCoupon(ctx context.Context, id string) (models.Coupon, error)
	ListCoupons(ctx context.Context, limit int32) ([]models.Coupon, error)
	UpdateCoupon(ctx context.Context, coupon models.Coupon) (models.Coupon, error)
	DeleteCoupon(ctx context.Context, id string) error
}

// couponsListLimit is how many a listing returns, newest first. Paging is
// yours to add when a client needs more.
const couponsListLimit = 100

// CouponsController is the HTTP surface of coupons: a JSON resource over the store. A
// handler converts HTTP into models, calls one store method, and
// converts the result back. Validation is the model's and runs on save,
// which is why a bad value is a 422 from the store and not a binding rule
// here. Nothing below app/controllers sees *gin.Context.
type CouponsController struct {
	Store couponsStore
	Log   *slog.Logger
}

// NewCouponsController builds the controller.
func NewCouponsController(store couponsStore, log *slog.Logger) *CouponsController {
	return &CouponsController{Store: store, Log: log}
}

// SetupRoutes mounts the resource under its own prefix. Middleware the whole
// controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by routes.go, which calls this.
func (ctl *CouponsController) SetupRoutes(router *gin.RouterGroup) {
	coupons := router.Group("/coupons")
	coupons.GET("", ctl.Index)
	coupons.POST("", ctl.Create)
	coupons.GET("/:id", ctl.Show)
	coupons.PATCH("/:id", ctl.Update)
	coupons.DELETE("/:id", ctl.Destroy)
}

// Index handles GET /admin/coupons.
func (ctl *CouponsController) Index(c *gin.Context) {
	coupons, err := ctl.Store.ListCoupons(c.Request.Context(), couponsListLimit)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, couponsView(coupons))
}

// Show handles GET /admin/coupons/:id.
func (ctl *CouponsController) Show(c *gin.Context) {
	coupon, err := ctl.Store.GetCoupon(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, couponView(coupon))
}

// Create handles POST /admin/coupons.
func (ctl *CouponsController) Create(c *gin.Context) {
	var req views.AdminCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.Error(err))
		return
	}
	coupon, err := ctl.Store.CreateCoupon(c.Request.Context(), models.Coupon{Code: req.Code, PriceRuleID: req.PriceRuleID})
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, couponView(coupon))
}

// Update handles PATCH /admin/coupons/:id. Only the fields in the body change: the coupon is loaded,
// the body is decoded over its current values,
// so a key left out keeps its value, and the result is saved and validated
// whole.
func (ctl *CouponsController) Update(c *gin.Context) {
	coupon, err := ctl.Store.GetCoupon(c.Request.Context(), c.Param("id"))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	req := views.AdminCouponRequest{Code: coupon.Code, PriceRuleID: coupon.PriceRuleID}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.Error(err))
		return
	}
	coupon, err = ctl.Store.UpdateCoupon(c.Request.Context(), models.Coupon{ID: coupon.ID, Code: req.Code, PriceRuleID: req.PriceRuleID})
	if err != nil {
		ctl.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, couponView(coupon))
}

// Destroy handles DELETE /admin/coupons/:id.
func (ctl *CouponsController) Destroy(c *gin.Context) {
	if err := ctl.Store.DeleteCoupon(c.Request.Context(), c.Param("id")); err != nil {
		ctl.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// fail maps a model error to a status, the one place that mapping lives:
// not found is 404, a validation failure is 422, and anything
// else is logged and answered 500 without leaking its text.
func (ctl *CouponsController) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		c.JSON(http.StatusNotFound, views.ErrorMessage("no such coupon"))
	case errors.Is(err, models.ErrInvalid):
		c.JSON(http.StatusUnprocessableEntity, views.Error(err))
	default:
		ctl.Log.ErrorContext(c.Request.Context(), "coupons: request failed",
			"path", c.FullPath(), "err", err, "request_id", c.GetString("request_id"))
		c.JSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
	}
}

// The controller converts between the models and the wire. views never
// imports models, so this is where AdminCouponResponse is built.

func couponView(coupon models.Coupon) views.AdminCouponResponse {
	return views.AdminCouponResponse{
		ID:          coupon.ID,
		Code:        coupon.Code,
		PriceRuleID: coupon.PriceRuleID,
		CreatedAt:   coupon.CreatedAt,
		UpdatedAt:   coupon.UpdatedAt,
	}
}

func couponsView(coupons []models.Coupon) views.AdminCouponsResponse {
	out := views.AdminCouponsResponse{Coupons: make([]views.AdminCouponResponse, 0, len(coupons))}
	for _, coupon := range coupons {
		out.Coupons = append(out.Coupons, couponView(coupon))
	}
	return out
}
