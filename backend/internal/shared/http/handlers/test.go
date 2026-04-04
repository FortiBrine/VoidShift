package handlers

import (
	"net/http"

	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/labstack/echo/v5"
)

func TestHandler(c *echo.Context) error {
	u, ok := c.Get(auth.ContextUserIDKey).(*user.User)

	if !ok || u == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"user_id": u.ID,
	})
}
