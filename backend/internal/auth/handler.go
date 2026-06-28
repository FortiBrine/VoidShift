package auth

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

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

func (h *LoginHandler) Login(c fiber.Ctx) error {
	req := new(LoginRequestDto)
	if err := c.Bind().JSON(req); err != nil {
		return err
	}

	sess := session.FromContext(c)
	if sess.Get("username") != nil {
		return c.SendStatus(http.StatusConflict)
	}

	err := h.authService.Login(
		c.Context(),
		req.Username,
		req.Password,
		c.UserAgent(),
		c.IP(),
	)

	if errors.Is(err, ErrWrongPassword) || errors.Is(err, ErrUserNotFound) {
		return c.SendStatus(http.StatusUnauthorized)
	}
	if err != nil {
		return err
	}

	sess.Set("username", req.Username)

	return c.SendStatus(http.StatusNoContent)
}

func (h *LoginHandler) Logout(c fiber.Ctx) error {
	sess := session.FromContext(c)
	if sess.Get("username") == nil {
		return c.SendStatus(http.StatusUnauthorized)
	}
	sess.Delete("username")

	return c.SendStatus(http.StatusNoContent)
}
