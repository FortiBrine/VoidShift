package webui

import (
	"errors"

	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/FortiBrine/VoidShift/internal/validator"
	"github.com/FortiBrine/VoidShift/internal/webui"
	"github.com/FortiBrine/VoidShift/view/pages"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type Handler struct {
	authService *auth.Service
	i18nService *i18n.Service
}

func NewHandler(
	authService *auth.Service,
	i18nService *i18n.Service,
) *Handler {
	return &Handler{
		authService: authService,
		i18nService: i18nService,
	}
}

type LoginForm struct {
	Username string `form:"username" validate:"required,min=4,max=30,alphanum" message:"login.validation.username"`
	Password string `form:"password" validate:"required,min=8,max=40" message:"login.validation.password"`
}

func (h *Handler) LoginPage(c fiber.Ctx) error {
	if webui.IsAuthenticated(c) {
		return c.Redirect().To("/")
	}
	return webui.Render(c, pages.Login(webui.Localizer(h.i18nService, c), "", "", nil))
}

func (h *Handler) Login(c fiber.Ctx) error {
	if webui.IsAuthenticated(c) {
		return c.Redirect().To("/")
	}

	loc := webui.Localizer(h.i18nService, c)

	req := new(LoginForm)
	if err := c.Bind().Form(req); err != nil {
		if validationErr, ok := errors.AsType[*validator.ValidationError](err); ok {
			return webui.Render(c, pages.Login(loc, req.Username, "", validationErr.Fields))
		}
		return err
	}

	sess := session.FromContext(c)

	err := h.authService.Login(c.Context(), req.Username, req.Password, c.UserAgent(), c.IP())
	if errors.Is(err, auth.ErrWrongPassword) || errors.Is(err, auth.ErrUserNotFound) {
		return webui.Render(c, pages.Login(loc, req.Username, i18n.T(loc, "login.error_invalid_credentials"), nil))
	}
	if err != nil {
		return err
	}

	sess.Set("username", req.Username)

	return c.Redirect().To("/")
}

func (h *Handler) Logout(c fiber.Ctx) error {
	session.FromContext(c).Delete("username")
	return c.Redirect().To("/login")
}
