package auth

import (
	"errors"
	"net/http"

	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/labstack/echo/v5"
)

const ContextUserIDKey = "userID"

func Middleware(
	sessionService *session.Service,
	userService *user.Service,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			cookie, err := c.Cookie(SessionCookieName)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "missing session cookie")
			}

			userID, err := sessionService.ValidateSession(c.Request().Context(), cookie.Value)
			if err != nil {
				if errors.Is(err, session.ErrInvalidSession) {
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid session")
				}

				c.Echo().Logger.Error("failed to validate session", "error", err)
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
