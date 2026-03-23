package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FortiBrine/VoidShift/internal/http/handlers/auth"
	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/labstack/echo/v5"
)

const ContextUserIDKey = "userID"

func AuthMiddleware(
	sessionService *session.Service,
	userService *user.Service,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			cookie, err := c.Cookie(auth.SessionCookieName)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "missing session cookie")
			}

			userID, err := sessionService.ValidateSession(c.Request().Context(), cookie.Value)
			if err != nil {
				if errors.Is(err, session.ErrInvalidSession) {
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid session")
				}

				fmt.Printf("failed to validate session: %v\n", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to validate session")
			}

			u, err := userService.GetByID(c.Request().Context(), userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
			}

			c.Set(ContextUserIDKey, u)

			return next(c)
		}
	}
}
