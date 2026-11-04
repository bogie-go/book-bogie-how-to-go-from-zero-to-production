// Package admin_controllers is the admin namespace: the controllers mounted under
// /admin. One file per
// controller; the group and its shared middleware in routes.go.
package admin_controllers

import (
	"log/slog"
)

// Server holds the namespace's controllers, one NAMED field each, so adding
// one is one line here, one in app/application.go and one in routes.go.
// "bogie g controller admin/NAME" inserts all three.
type Server struct {
	Log *slog.Logger

	Products   *ProductsController
	PriceRules *PriceRulesController
	Coupons    *CouponsController
	// bogie:controllers
}

// NewServer builds the namespace. Its controllers are constructed and
// assigned in app/application.go, where every dependency is wired.
func NewServer(log *slog.Logger) *Server {
	return &Server{Log: log}
}
