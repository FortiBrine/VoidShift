package wireguard

import (
	"fmt"
	"net/http"

	"github.com/FortiBrine/VoidShift/internal/shared"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

type GenerateNetworkRequest struct {
	Name       string `json:"name" validate:"required,min=4,max=100"`
	Address    string `json:"address" validate:"required,cidr"`
	ListenPort int    `json:"listen_port" validate:"required,min=1024,max=65535"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetNetworks(c *echo.Context) error {
	ctx := c.Request().Context()

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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"networks": networksResult,
	})
}

func (h *Handler) GenerateNetwork(c *echo.Context) error {
	ctx := c.Request().Context()
	var request GenerateNetworkRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	network, err := h.service.GenerateNetwork(ctx, request.Name, request.Address, request.ListenPort)
	if err != nil {
		return fmt.Errorf("failed to generate network: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":          network.ID,
		"public_key":  network.PublicKey,
		"address":     network.Address,
		"listen_port": network.ListenPort,
	})
}

func (h *Handler) GetNetwork(c *echo.Context) error {
	ctx := c.Request().Context()
	networkID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return shared.ErrNetworkNotFound
	}

	network, err := h.service.GetNetworkWithPeers(ctx, networkID)
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

	return c.JSON(http.StatusOK, map[string]any{
		"id":          network.ID,
		"public_key":  network.PublicKey,
		"address":     network.Address,
		"listen_port": network.ListenPort,

		"peers": peers,
	})

}

func (h *Handler) RemoveNetwork(c *echo.Context) error {
	ctx := c.Request().Context()
	networkID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return shared.ErrNetworkNotFound
	}

	if err := h.service.RemoveNetwork(ctx, networkID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) UpNetwork(c *echo.Context) error {
	ctx := c.Request().Context()
	networkID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return shared.ErrNetworkNotFound
	}

	if err := h.service.UpNetwork(ctx, networkID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DownNetwork(c *echo.Context) error {
	ctx := c.Request().Context()
	networkID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return shared.ErrNetworkNotFound
	}

	if err := h.service.DownNetwork(ctx, networkID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GeneratePeer(c *echo.Context) error {
	ctx := c.Request().Context()
	networkID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return shared.ErrNetworkNotFound
	}

	peer, err := h.service.GeneratePeer(ctx, networkID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":         peer.ID,
		"public_key": peer.PublicKey,
	})
}

func (h *Handler) RemovePeer(c *echo.Context) error {
	ctx := c.Request().Context()
	peerID, err := echo.PathParam[uint](c, "peerId")
	if err != nil {
		return shared.ErrPeerNotFound
	}

	if err := h.service.RemovePeer(ctx, peerID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
