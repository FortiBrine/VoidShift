package webui

import (
	"fmt"
	"time"

	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/FortiBrine/VoidShift/internal/webui"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/FortiBrine/VoidShift/view/pages"
	"github.com/gofiber/fiber/v3"
	"github.com/shirou/gopsutil/v4/host"
)

type Handler struct {
	wireGuardService *wireguard.Service
	i18nService      *i18n.Service
	hostAddress      string
}

func NewHandler(
	wireGuardService *wireguard.Service,
	i18nService *i18n.Service,
	hostAddress string,
) *Handler {
	return &Handler{
		wireGuardService: wireGuardService,
		i18nService:      i18nService,
		hostAddress:      hostAddress,
	}
}

func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}

func (h *Handler) Home(c fiber.Ctx) error {
	stats, err := h.wireGuardService.GetStats(c.Context())
	if err != nil {
		return fmt.Errorf("getting wireguard stats: %w", err)
	}

	uptimeSeconds, err := host.Uptime()
	if err != nil {
		return fmt.Errorf("getting uptime: %w", err)
	}

	uptime := time.Duration(uptimeSeconds) * time.Second
	return webui.Render(c, pages.Dashboard(webui.Localizer(h.i18nService, c), h.hostAddress, formatUptime(uptime), stats))
}
