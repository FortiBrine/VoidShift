package webui

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/FortiBrine/VoidShift/internal/validator"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/FortiBrine/VoidShift/view/pages"
	"github.com/gofiber/fiber/v3"
)

func redirectWithError(c fiber.Ctx, path string, message string) error {
	return c.Redirect().To(path + "?error=" + url.QueryEscape(message))
}

func (h *Handler) Networks(c fiber.Ctx) error {
	ctx := c.Context()

	networks, err := h.wireGuardService.GetNetworks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get networks: %w", err)
	}

	return Render(c, pages.Networks(networks, c.Query("error")))
}

type NetworkCreateForm struct {
	Name       string `form:"name" validate:"required,min=4,max=100" message:"Назва має бути від 4 до 100 символів"`
	Address    string `form:"address" validate:"required,cidr" message:"Адреса має бути в форматі CIDR"`
	ListenPort int    `form:"listen_port" validate:"required,min=1024,max=65535" message:"Порт має бути в межах 1024-65535"`
}

func (h *Handler) NetworkCreatePage(c fiber.Ctx) error {
	return Render(c, pages.NetworkCreate("", "10.8.0.1/24", 51820, nil))
}

func (h *Handler) NetworkCreate(c fiber.Ctx) error {
	req := new(NetworkCreateForm)
	if err := c.Bind().Form(req); err != nil {
		if validationErr, ok := errors.AsType[*validator.ValidationError](err); ok {
			return Render(c, pages.NetworkCreate(req.Name, req.Address, req.ListenPort, validationErr.Fields))
		}
		return err
	}

	ctx := c.Context()
	network, err := h.wireGuardService.GenerateNetwork(ctx, req.Name, req.Address, req.ListenPort)
	if err != nil {
		return fmt.Errorf("failed to generate network: %w", err)
	}

	return c.Redirect().To(fmt.Sprintf("/wireguard/networks/%d", network.ID))
}

func (h *Handler) NetworkDetail(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	network, err := h.wireGuardService.GetNetworkWithPeers(ctx, networkID)
	if errors.Is(err, wireguard.ErrNetworkNotFound) {
		return redirectWithError(c, "/wireguard", "Мережу не знайдено")
	}
	if err != nil {
		return fmt.Errorf("failed to get network: %w", err)
	}

	return Render(c, pages.NetworkDetail(network))
}

func (h *Handler) NetworkUp(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	if err := h.wireGuardService.UpNetwork(ctx, networkID); err != nil {
		if errors.Is(err, wireguard.ErrNetworkNotFound) {
			return redirectWithError(c, "/wireguard", "Мережу не знайдено")
		}
		return fmt.Errorf("failed to bring network up: %w", err)
	}

	return c.Redirect().To(fmt.Sprintf("/wireguard/networks/%d", networkID))
}

func (h *Handler) NetworkDown(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	if err := h.wireGuardService.DownNetwork(ctx, networkID); err != nil {
		if errors.Is(err, wireguard.ErrNetworkNotFound) {
			return redirectWithError(c, "/wireguard", "Мережу не знайдено")
		}
		return fmt.Errorf("failed to bring network down: %w", err)
	}

	return c.Redirect().To(fmt.Sprintf("/wireguard/networks/%d", networkID))
}

type PeerCreateForm struct {
	IP string `form:"ip" validate:"required,ipv4" message:"Вкажи коректну IPv4 адресу"`
}

func (h *Handler) PeerCreatePage(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	if _, err := h.wireGuardService.GetNetwork(ctx, networkID); err != nil {
		if errors.Is(err, wireguard.ErrNetworkNotFound) {
			return redirectWithError(c, "/wireguard", "Мережу не знайдено")
		}
		return fmt.Errorf("failed to get network: %w", err)
	}

	return Render(c, pages.PeerCreate(networkID, "", nil))
}

func (h *Handler) PeerCreate(c fiber.Ctx) error {
	networkID := fiber.Params[uint](c, "id")

	req := new(PeerCreateForm)
	if err := c.Bind().Form(req); err != nil {
		if validationErr, ok := errors.AsType[*validator.ValidationError](err); ok {
			return Render(c, pages.PeerCreate(networkID, req.IP, validationErr.Fields))
		}
		return err
	}

	ctx := c.Context()
	if _, err := h.wireGuardService.GeneratePeer(ctx, networkID, []string{req.IP}); err != nil {
		if errors.Is(err, wireguard.ErrNetworkNotFound) {
			return redirectWithError(c, "/wireguard", "Мережу не знайдено")
		}
		return fmt.Errorf("failed to generate peer: %w", err)
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
			return redirectWithError(c, "/wireguard", "Peer не знайдено")
		}
		return fmt.Errorf("failed to get peer config: %w", err)
	}

	return Render(c, pages.PeerConfig(peerID, networkID, config))
}

func (h *Handler) PeerQR(c fiber.Ctx) error {
	peerID := fiber.Params[uint](c, "peerId")
	networkID := fiber.Query[uint](c, "networkId", 0)

	return Render(c, pages.PeerQR(peerID, networkID))
}
