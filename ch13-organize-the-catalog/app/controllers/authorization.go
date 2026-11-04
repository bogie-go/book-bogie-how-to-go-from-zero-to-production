package controllers

import (
	"github.com/gin-gonic/gin"

	"depot/app/middlewares"
)

// WithStaff is the route group only staff may use: a signed-in user, who is
// also on the shop's staff.
func (s *Server) WithStaff(r *gin.RouterGroup) *gin.RouterGroup {
	return s.WithUser(r).Group("", middlewares.RequireStaff(s.Staff, s.Log))
}
