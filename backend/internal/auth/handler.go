package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequestDto struct {
	Username string `json:"username" validate:"required,min=4,max=30"`
	Password string `json:"password" validate:"required,min=8,max=40"`
}

type LoginHandler struct {
	sessionService *session.Service
	userService    *user.Service
}

func NewLoginHandler(
	sessionService *session.Service,
	userService *user.Service,
) *LoginHandler {
	return &LoginHandler{
		sessionService: sessionService,
		userService:    userService,
	}
}

func (h *LoginHandler) Login(c *echo.Context) error {
	req := new(LoginRequestDto)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	u, err := h.userService.GetByUsername(ctx, req.Username)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	sessionID, expiresAt, err := h.sessionService.CreateUserSession(
		ctx,
		u.ID, c.Request().UserAgent(), c.RealIP(),
	)

	if err != nil {
		fmt.Printf("error creating session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	c.SetCookie(BuildSessionCookie(sessionID, expiresAt))

	return c.JSON(http.StatusOK, map[string]string{
		"message": "success",
	})

}

func (h *LoginHandler) Logout(c *echo.Context) error {
	ctx := c.Request().Context()
	cookie, err := c.Cookie(SessionCookieName)
	if err == nil {
		err := h.sessionService.LogoutSession(ctx, cookie.Value)

		if err != nil {
			fmt.Printf("error logout: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
	}

	c.SetCookie(BuildExpiredSessionCookie())
	return c.String(http.StatusNoContent, "")

}
