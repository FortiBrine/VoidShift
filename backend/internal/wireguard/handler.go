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

func (h *Handler) GetNetwork(c *echo.Context) error {
	ctx := c.Request().Context()

	networks, err := h.service.GetNetworks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get networks: %w", err)
	}

	result := make(map[uint]interface{})
	for _, network := range networks {
		result[network.ID] = map[string]interface{}{
			"Name":       network.Name,
			"Address":    network.Address,
			"ListenPort": network.ListenPort,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"networks": result,
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
		"network": network.ID,
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
