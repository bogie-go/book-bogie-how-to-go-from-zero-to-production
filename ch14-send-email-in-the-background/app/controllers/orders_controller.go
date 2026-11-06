package controllers

import (
	"cmp"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"depot/app/jobs"
	"depot/app/middlewares"
	"depot/app/models"
	"depot/app/views"
	"depot/db/adapter"
)

type ordersStore interface {
	FindOrCreateCartByToken(ctx context.Context, token string) (models.Cart, error)
	ListActivePriceRules(ctx context.Context) ([]models.PriceRule, error)
	CouponRule(ctx context.Context, code string) (models.PriceRule, error)
	FindOrderForUser(ctx context.Context, id string, userID int64) (models.Order, error)
	PlaceOrder(ctx context.Context, userID, cartID int64, rules []models.PriceRule, after adapter.AfterOrder) (models.Order, error)
	OrderLines(ctx context.Context, orderID int64) ([]models.OrderItem, error)
}

// jobQueue is the slice of the River client this controller uses.
type jobQueue interface {
	InsertTx(ctx context.Context, tx pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// OrdersController is the HTTP surface of orders: its dependencies, its routes
// and its actions, all in one file.
type OrdersController struct {
	Store ordersStore
	Jobs  jobQueue
	Log   *slog.Logger
}

// NewOrdersController builds the controller.
func NewOrdersController(store ordersStore, jobs jobQueue, log *slog.Logger) *OrdersController {
	return &OrdersController{Store: store, Jobs: jobs, Log: log}
}

// SetupRoutes mounts the controller under its own prefix. Middleware the
// whole controller needs goes here; which group it sits in, and so how it is
// exposed, is decided by the routes.go that calls this.
func (ctl *OrdersController) SetupRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	orders.GET("/:id", ctl.Show)
	orders.POST("", ctl.Create)
}

// Show handles GET /orders/:id: one of the customer's own orders.
//
// A handler converts HTTP into models, calls one service function,
// and converts the result back. Nothing below app/controllers sees
// *gin.Context.
func (ctl *OrdersController) Show(c *gin.Context) {
	userID, ok := models.UserIDFromSubject(middlewares.Subject(c))
	if !ok {
		ctl.fail(c, models.ErrNotFound)
		return
	}
	order, err := ctl.Store.FindOrderForUser(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	ctl.respondWithOrder(c, http.StatusOK, order)
}

// Create handles POST /orders: check out the customer's cart.
func (ctl *OrdersController) Create(c *gin.Context) {
	userID, ok := models.UserIDFromSubject(middlewares.Subject(c))
	if !ok {
		c.JSON(http.StatusForbidden, views.ErrorMessage("create an account to place an order"))
		return
	}
	var req views.PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorMessage("invalid JSON"))
		return
	}
	ctx := c.Request.Context()
	cart, err := ctl.Store.FindOrCreateCartByToken(ctx, middlewares.Subject(c))
	if err != nil {
		ctl.fail(c, err)
		return
	}
	rules, err := ctl.Store.ListActivePriceRules(ctx)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	if code := models.NormalizeCouponCode(req.CouponCode); code != "" {
		rule, err := ctl.Store.CouponRule(ctx, code)
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusUnprocessableEntity, views.ErrorMessage("that coupon code isn't valid"))
			return
		}
		if err != nil {
			ctl.fail(c, err)
			return
		}
		rules = append(rules, rule)
		slices.SortStableFunc(rules, func(a, b models.PriceRule) int { return cmp.Compare(a.Priority, b.Priority) })
	}
	order, err := ctl.Store.PlaceOrder(ctx, userID, cart.ID, rules,
		func(ctx context.Context, tx pgx.Tx, order models.Order) error {
			_, err := ctl.Jobs.InsertTx(ctx, tx, jobs.SendOrderConfirmationArgs{OrderID: order.ID}, nil)
			return err
		})
	if err != nil {
		ctl.fail(c, err)
		return
	}
	ctl.respondWithOrder(c, http.StatusCreated, order)
}

// respondWithOrder answers with an order and its lines.
func (ctl *OrdersController) respondWithOrder(c *gin.Context, status int, order models.Order) {
	items, err := ctl.Store.OrderLines(c.Request.Context(), order.ID)
	if err != nil {
		ctl.fail(c, err)
		return
	}
	resp := views.OrderResponse{
		ID: order.ID, Status: order.Status, SubtotalCents: order.SubtotalCents,
		DiscountCents: order.DiscountCents, TotalCents: order.TotalCents,
		Items: make([]views.OrderItemResponse, 0, len(items)), CreatedAt: order.CreatedAt,
	}
	for _, i := range items {
		resp.Items = append(resp.Items, views.OrderItemResponse{
			ProductID: i.ProductID, Name: i.Name, Quantity: i.Quantity,
			UnitPriceCents: i.UnitPriceCents, SubtotalCents: i.UnitPriceCents * i.Quantity,
		})
	}
	c.JSON(status, resp)
}

// fail turns an order error into the right status.
func (ctl *OrdersController) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		c.JSON(http.StatusNotFound, views.ErrorMessage("no such order"))
	case errors.Is(err, models.ErrInvalid):
		c.JSON(http.StatusUnprocessableEntity, views.Error(err))
	default:
		ctl.Log.ErrorContext(c.Request.Context(), "orders: request failed", "path", c.FullPath(), "err", err, "request_id", c.GetString("request_id"))
		c.JSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
	}
}
