package wireguard

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
}

type GenerateNetworkRequest struct {
	Name       string `json:"name" validate:"required,min=4,max=100" message:"Name must be between 4 and 100 characters"`
	Address    string `json:"address" validate:"required,cidr" message:"Address must be CIDR"`
	ListenPort int    `json:"listen_port" validate:"required,min=1024,max=65535" message:"ListenPort must be between 1024 and 65535"`
}

type GeneratePeerRequest struct {
	AllowedIPs []string `json:"allowed_ips" validate:"required,dive,ipv4" message:"AllowedIPs must be valid IPv4 addresses"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetNetworks(c fiber.Ctx) error {
	ctx := c.Context()

	networks, err := h.service.GetNetworks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get networks: %w", err)
	}

	networksResult := make([]map[string]any, len(networks))
	for i, network := range networks {
		networksResult[i] = map[string]any{
			"id":          network.ID,
			"name":        network.Name,
			"address":     network.Address,
			"listen_port": network.ListenPort,
		}
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"networks": networksResult,
	})
}

func (h *Handler) GenerateNetwork(c fiber.Ctx) error {
	ctx := c.Context()
	request := new(GenerateNetworkRequest)
	if err := c.Bind().JSON(request); err != nil {
		return err
	}

	network, err := h.service.GenerateNetwork(ctx, request.Name, request.Address, request.ListenPort)
	if err != nil {
		return fmt.Errorf("failed to generate network: %w", err)
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]any{
		"id":          network.ID,
		"public_key":  network.PublicKey,
		"address":     network.Address,
		"listen_port": network.ListenPort,
	})
}

func (h *Handler) GetNetwork(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	network, err := h.service.GetNetworkWithPeers(ctx, networkID)
	if errors.Is(err, ErrNetworkNotFound) {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		return err
	}

	peers := make([]map[string]any, len(network.Peers))
	for i, peer := range network.Peers {
		peers[i] = map[string]any{
			"id":          peer.ID,
			"public_key":  peer.PublicKey,
			"allowed_ips": peer.AllowedIPs,
		}
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"id":          network.ID,
		"public_key":  network.PublicKey,
		"address":     network.Address,
		"listen_port": network.ListenPort,

		"peers": peers,
	})

}

func (h *Handler) RemoveNetwork(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	if err := h.service.RemoveNetwork(ctx, networkID); err != nil {
		if errors.Is(err, ErrNetworkNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) UpNetwork(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	if err := h.service.UpNetwork(ctx, networkID); err != nil {
		if errors.Is(err, ErrNetworkNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DownNetwork(c fiber.Ctx) error {
	ctx := c.Context()
	networkID := fiber.Params[uint](c, "id")

	if err := h.service.DownNetwork(ctx, networkID); err != nil {
		if errors.Is(err, ErrNetworkNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GeneratePeer(c fiber.Ctx) error {
	ctx := c.Context()
	request := new(GeneratePeerRequest)
	if err := c.Bind().JSON(request); err != nil {
		return err
	}

	networkID := fiber.Params[uint](c, "id")

	peer, err := h.service.GeneratePeer(ctx, networkID, request.AllowedIPs)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]any{
		"id":         peer.ID,
		"public_key": peer.PublicKey,
	})
}

func (h *Handler) RemovePeer(c fiber.Ctx) error {
	ctx := c.Context()
	peerID := fiber.Params[uint](c, "peerId")

	if err := h.service.RemovePeer(ctx, peerID); err != nil {
		if errors.Is(err, ErrPeerNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetPeerConfig(c fiber.Ctx) error {
	ctx := c.Context()
	peerID := fiber.Params[uint](c, "peerId")

	config, err := h.service.GetPeerConfig(ctx, peerID)
	if err != nil {
		if errors.Is(err, ErrPeerNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.Status(fiber.StatusOK).SendString(config)
}

func (h *Handler) DownloadPeerConfig(c fiber.Ctx) error {
	ctx := c.Context()
	peerID := fiber.Params[uint](c, "peerId")

	config, err := h.service.GetPeerConfig(ctx, peerID)
	if err != nil {
		if errors.Is(err, ErrPeerNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	c.Response().Header.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"peer-%d.conf\"", peerID))
	return c.Status(fiber.StatusOK).SendString(config)
}

func (h *Handler) GetPeerConfigQR(c fiber.Ctx) error {
	ctx := c.Context()
	peerID := fiber.Params[uint](c, "peerId")

	qrCode, err := h.service.GetPeerConfigQR(ctx, peerID)
	if err != nil {
		if errors.Is(err, ErrPeerNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.Status(fiber.StatusOK).
		Type("png").
		Send(qrCode)
}
