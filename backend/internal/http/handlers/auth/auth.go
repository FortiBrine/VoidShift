package auth

import (
	"net/http"

	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
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
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	user, err := h.userService.GetByUsername(ctx, req.Username)

	if err != nil {
		return c.String(http.StatusNotFound, "user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))

	if err != nil {
		return c.String(http.StatusNotFound, "user not found")
	}

	sessionID, expiresAt, err := h.sessionService.CreateUserSession(
		ctx,
		user.ID, c.Request().UserAgent(), c.RealIP(),
	)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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
			return c.String(http.StatusInternalServerError, "internal server error")
		}
	}

	c.SetCookie(BuildExpiredSessionCookie())
	return c.String(http.StatusNoContent, "")

}
