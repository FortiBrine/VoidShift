package webui

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/FortiBrine/VoidShift/internal/validator"
	"github.com/FortiBrine/VoidShift/internal/webui"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	pages "github.com/FortiBrine/VoidShift/view/pages/wireguard"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	wireGuardService *wireguard.Service
	i18nService      *i18n.Service
}

func NewHandler(
	wireGuardService *wireguard.Service,
	i18nService *i18n.Service,
) *Handler {
	return &Handler{
		wireGuardService: wireGuardService,
		i18nService:      i18nService,
	}
}

func (h *Handler) redirectWithError(c fiber.Ctx, path string, messageID string) error {
	message := i18n.T(webui.Localizer(h.i18nService, c), messageID)
	return c.Redirect().To(path + "?error=" + url.QueryEscape(message))
}

func (h *Handler) Networks(c fiber.Ctx) error {
	ctx := c.Context()

	networks, err := h.wireGuardService.GetNetworks(ctx)
	if err != nil {
		return fmt.Errorf("getting networks: %w", err)
	}

	return webui.Render(c, pages.Networks(webui.Localizer(h.i18nService, c), networks, c.Query("error")))
}

type NetworkCreateForm struct {
	Name       string `form:"name" validate:"required,min=4,max=100" message:"network_create.validation.name"`
	Address    string `form:"address" validate:"required,cidr" message:"network_create.validation.address"`
	ListenPort int    `form:"listen_port" validate:"required,min=1024,max=65535" message:"network_create.validation.listen_port"`
}

func (h *Handler) NetworkCreatePage(c fiber.Ctx) error {
	return webui.Render(c, pages.NetworkCreate(webui.Localizer(h.i18nService, c), "", "10.8.0.1/24", 51820, nil))
}

func (h *Handler) NetworkCreate(c fiber.Ctx) error {
	req := new(NetworkCreateForm)
	if err := c.Bind().Form(req); err != nil {
		if validationErr, ok := errors.AsType[*validator.ValidationError](err); ok {
			return webui.Render(c, pages.NetworkCreate(webui.Localizer(h.i18nService, c), req.Name, req.Address, req.ListenPort, validationErr.Fields))
		}
		return err
	}

	ctx := c.Context()
	network, err := h.wireGuardService.GenerateNetwork(ctx, req.Name, req.Address, req.ListenPort)
	if err != nil {
		return fmt.Errorf("generating network: %w", err)
	}

	return c.Redirect().To(fmt.Sprintf("/wireguard/networks/%d", network.ID))
}

func (h *Handler) NetworkDetail(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	network, err := h.wireGuardService.GetNetworkWithPeers(ctx, networkID)
	if errors.Is(err, wireguard.ErrNetworkNotFound) {
		return h.redirectWithError(c, "/wireguard", "wireguard.network_not_found")
	}
	if err != nil {
		return fmt.Errorf("getting network: %w", err)
	}

	return webui.Render(c, pages.NetworkDetail(webui.Localizer(h.i18nService, c), network))
}

func (h *Handler) NetworkUp(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	if err := h.wireGuardService.UpNetwork(ctx, networkID); err != nil {
		if errors.Is(err, wireguard.ErrNetworkNotFound) {
			return h.redirectWithError(c, "/wireguard", "wireguard.network_not_found")
		}
		return fmt.Errorf("bringing network up: %w", err)
	}

	return c.Redirect().To(fmt.Sprintf("/wireguard/networks/%d", networkID))
}

func (h *Handler) NetworkDown(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	if err := h.wireGuardService.DownNetwork(ctx, networkID); err != nil {
		if errors.Is(err, wireguard.ErrNetworkNotFound) {
			return h.redirectWithError(c, "/wireguard", "wireguard.network_not_found")
		}
		return fmt.Errorf("bringing network down: %w", err)
	}

	return c.Redirect().To(fmt.Sprintf("/wireguard/networks/%d", networkID))
}

type PeerCreateForm struct {
	IP string `form:"ip" validate:"required,ipv4" message:"peer_create.validation.ip"`
}

func (h *Handler) PeerCreatePage(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	if _, err := h.wireGuardService.GetNetwork(ctx, networkID); err != nil {
		if errors.Is(err, wireguard.ErrNetworkNotFound) {
			return h.redirectWithError(c, "/wireguard", "wireguard.network_not_found")
		}
		return fmt.Errorf("getting network: %w", err)
	}

	return webui.Render(c, pages.PeerCreate(webui.Localizer(h.i18nService, c), networkID, "", nil))
}

func (h *Handler) PeerCreate(c fiber.Ctx) error {
	networkID := fiber.Params[uint](c, "id")

	req := new(PeerCreateForm)
	if err := c.Bind().Form(req); err != nil {
		if validationErr, ok := errors.AsType[*validator.ValidationError](err); ok {
			return webui.Render(c, pages.PeerCreate(webui.Localizer(h.i18nService, c), networkID, req.IP, validationErr.Fields))
		}
		return err
	}

	ctx := c.Context()
	if _, err := h.wireGuardService.GeneratePeer(ctx, networkID, []string{req.IP}); err != nil {
		if errors.Is(err, wireguard.ErrNetworkNotFound) {
			return h.redirectWithError(c, "/wireguard", "wireguard.network_not_found")
		}
		return fmt.Errorf("generating peer: %w", err)
	}

	return c.Redirect().To(fmt.Sprintf("/wireguard/networks/%d", networkID))
}

func (h *Handler) PeerConfig(c fiber.Ctx) error {
	ctx := c.Context()
	peerID := fiber.Params[uint](c, "peerId")
	networkID := fiber.Query[uint](c, "networkId", 0)

	config, err := h.wireGuardService.GetPeerConfig(ctx, peerID)
	if err != nil {
		if errors.Is(err, wireguard.ErrPeerNotFound) {
			return h.redirectWithError(c, "/wireguard", "wireguard.peer_not_found")
		}
		return fmt.Errorf("getting peer config: %w", err)
	}

	return webui.Render(c, pages.PeerConfig(webui.Localizer(h.i18nService, c), peerID, networkID, config))
}

func (h *Handler) PeerQR(c fiber.Ctx) error {
	peerID := fiber.Params[uint](c, "peerId")
	networkID := fiber.Query[uint](c, "networkId", 0)

	return webui.Render(c, pages.PeerQR(webui.Localizer(h.i18nService, c), peerID, networkID))
}
