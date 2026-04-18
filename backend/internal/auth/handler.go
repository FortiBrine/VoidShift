package auth

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

const SessionCookieName = "voidshift.session"

type LoginRequestDto struct {
	Username string `json:"username" validate:"required,min=4,max=30,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=40"`
}

type LoginHandler struct {
	authService *Service
}

func NewLoginHandler(
	authService *Service,
) *LoginHandler {
	return &LoginHandler{
		authService: authService,
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

	result, err := h.authService.Login(
		c.Request().Context(),
		req.Username,
		req.Password,
		c.Request().UserAgent(),
		c.RealIP(),
	)

	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     SessionCookieName,
		Value:    result.SessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  result.ExpiresAt,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(result.ExpiresAt).Seconds()),
	})

	return c.NoContent(http.StatusNoContent)
}

func (h *LoginHandler) Logout(c *echo.Context) error {
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "session not found")
	}

	if err := h.authService.LogoutSession(c.Request().Context(), cookie.Value); err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	return c.NoContent(http.StatusNoContent)
}
