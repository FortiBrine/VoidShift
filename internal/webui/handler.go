package webui

import (
	"errors"

	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/FortiBrine/VoidShift/internal/validator"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/FortiBrine/VoidShift/view/pages"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type Handler struct {
	authService      *auth.Service
	wireGuardService *wireguard.Service
}

func NewHandler(
	authService *auth.Service,
	wireGuardService *wireguard.Service,
) *Handler {
	return &Handler{
		authService:      authService,
		wireGuardService: wireGuardService,
	}
}

func isAuthenticated(c fiber.Ctx) bool {
	return session.FromContext(c).Get("username") != nil
}

func (h *Handler) Home(c fiber.Ctx) error {
	if isAuthenticated(c) {
		return c.Redirect().To("/wireguard")
	}
	return c.Redirect().To("/login")
}

type LoginForm struct {
	Username string `form:"username" validate:"required,min=4,max=30,alphanum" message:"login.validation.username"`
	Password string `form:"password" validate:"required,min=8,max=40" message:"login.validation.password"`
}

func (h *Handler) LoginPage(c fiber.Ctx) error {
	if isAuthenticated(c) {
		return c.Redirect().To("/wireguard")
	}
	return Render(c, pages.Login(localizer(c), "", "", nil))
}

func (h *Handler) Login(c fiber.Ctx) error {
	if isAuthenticated(c) {
		return c.Redirect().To("/wireguard")
	}

	loc := localizer(c)

	req := new(LoginForm)
	if err := c.Bind().Form(req); err != nil {
		if validationErr, ok := errors.AsType[*validator.ValidationError](err); ok {
			return Render(c, pages.Login(loc, req.Username, "", validationErr.Fields))
		}
		return err
	}

	sess := session.FromContext(c)

	err := h.authService.Login(c.Context(), req.Username, req.Password, c.UserAgent(), c.IP())
	if errors.Is(err, auth.ErrWrongPassword) || errors.Is(err, auth.ErrUserNotFound) {
		return Render(c, pages.Login(loc, req.Username, i18n.T(loc, "login.error_invalid_credentials"), nil))
	}
	if err != nil {
		return err
	}

	sess.Set("username", req.Username)

	return c.Redirect().To("/wireguard")
}

func (h *Handler) Logout(c fiber.Ctx) error {
	session.FromContext(c).Delete("username")
	return c.Redirect().To("/login")
}
