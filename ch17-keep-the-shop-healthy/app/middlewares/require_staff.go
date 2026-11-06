package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"depot/app/models"
	"depot/app/views"
)

// StaffLookup is what RequireStaff asks the store.
type StaffLookup interface {
	IsStaff(ctx context.Context, userID int64) (bool, error)
}

// RequireStaff lets a request through only for a member of staff. It runs
// after RequireUser, so the token is already verified; this asks the database
// whether its customer is staff, on every request, so taking someone's role
// away takes effect on their very next request.
func RequireStaff(staff StaffLookup, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if userID, ok := models.UserIDFromSubject(Subject(c)); ok {
			isStaff, err := staff.IsStaff(c.Request.Context(), userID)
			if err != nil {
				log.ErrorContext(c.Request.Context(), "require staff: lookup failed", "err", err, "request_id", c.GetString("request_id"))
				c.AbortWithStatusJSON(http.StatusInternalServerError, views.ErrorMessage("internal error"))
				return
			}
			if isStaff {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, views.ErrorMessage("staff only"))
	}
}
