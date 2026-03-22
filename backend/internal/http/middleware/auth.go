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
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}

			fmt.Printf("cookie value: %s\n", cookie.Value)

			userID, err := sessionService.ValidateSession(c.Request().Context(), cookie.Value)
			if err != nil {
				if errors.Is(err, session.ErrInvalidSession) {
					return c.JSON(http.StatusUnauthorized, "invalid session")
				}

				fmt.Printf("failed to validate session: %v\n", err)
				return c.JSON(http.StatusInternalServerError, "internal server error")
			}

			user, err := userService.GetByID(c.Request().Context(), userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, "internal server error")
			}

			c.Set(ContextUserIDKey, user)

			return next(c)
		}
	}
}
