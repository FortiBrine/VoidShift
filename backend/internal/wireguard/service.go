package wireguard

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/FortiBrine/VoidShift/internal/shared"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gorm.io/gorm"
)

type Service struct {
	repository Repository
	client     *wgctrl.Client
}

func NewService(repository Repository, client *wgctrl.Client) *Service {
	return &Service{
		repository: repository,
		client:     client,
	}
}

func (s *Service) Load() error {
	return s.repository.Migrate()
}

func (s *Service) GetNetworks(
	ctx context.Context,
) ([]Network, error) {
	return s.repository.GetNetworks(ctx)
}

func (s *Service) GenerateNetwork(
	ctx context.Context,
	name string,
	address string,
	listenPort int,
) (*Network, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey()
	network := &Network{
		Name: name,

		Address:    address,
		ListenPort: listenPort,

		PrivateKey: privateKey.String(),
		PublicKey:  publicKey.String(),
	}

	return network, s.repository.AddNetwork(ctx, network)
}

func (s *Service) GeneratePeer(
	ctx context.Context,
	networkID uint,
) (*Peer, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey()
	psk, err := wgtypes.GenerateKey()

	if err != nil {
		return nil, err
	}

	peer := &Peer{
		NetworkID: networkID,

		PublicKey:  publicKey.String(),
		PrivateKey: privateKey.String(),

		PresharedKey: psk.String(),
	}

	return peer, s.repository.AddPeer(ctx, peer)
}

func (s *Service) RemovePeer(
	ctx context.Context,
	peerID uint,
) error {
	_, err := s.repository.RemovePeer(ctx, peerID)
	return err
}

func (s *Service) RemoveNetwork(
	ctx context.Context,
	networkID uint,
) error {
	network, err := s.repository.GetNetwork(ctx, networkID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared.ErrNetworkNotFound
		}

		return fmt.Errorf("failed to get network: %w", err)
	}
	if err = RemoveDevice(network.Name); err != nil {
		return fmt.Errorf("failed to remove device: %w", err)
	}

	_, err = s.repository.RemoveNetwork(ctx, networkID)
	return err
}

func (s *Service) GetNetwork(
	ctx context.Context,
	networkID uint,
) (Network, error) {
	return s.repository.GetNetwork(ctx, networkID)
}

func (s *Service) UpNetwork(
	ctx context.Context,
	networkID uint,
) error {
	network, err := s.repository.GetNetworkWithPeers(ctx, networkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared.ErrNetworkNotFound
		}

		return err
	}

	if err := CreateDevice(network.Name); err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}

	if err := SetDeviceAddress(network.Name, network.Address); err != nil {
		return fmt.Errorf("failed to set device address: %w", err)
	}

	privateKey, err := wgtypes.ParseKey(network.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to parse private network key: %w", err)
	}

	peers := make([]wgtypes.PeerConfig, len(network.Peers))
	for i, peer := range network.Peers {
		publicKey, err := wgtypes.ParseKey(peer.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to parse public peer key: %w", err)
		}

		var presharedKey *wgtypes.Key
		if peer.PresharedKey != "" {
			psk, err := wgtypes.ParseKey(peer.PresharedKey)
			if err != nil {
				return fmt.Errorf("failed to parse preshared key: %w", err)
			}
			presharedKey = &psk
		}

		allowedIPs := make([]net.IPNet, len(peer.AllowedIPs))

		for j, allowedIP := range peer.AllowedIPs {
			_, ipNet, err := net.ParseCIDR(allowedIP)
			if err != nil {
				return fmt.Errorf("failed to parse allowed IP %q: %w", allowedIP, err)
			}

			allowedIPs[j] = *ipNet
		}

		peers[i] = wgtypes.PeerConfig{
			PublicKey:    publicKey,
			PresharedKey: presharedKey,
			AllowedIPs:   allowedIPs,
		}
	}

	if err := s.client.ConfigureDevice(network.Name, wgtypes.Config{
		PrivateKey:   &privateKey,
		ListenPort:   &network.ListenPort,
		ReplacePeers: true,
		Peers:        peers,
	}); err != nil {
		return fmt.Errorf("failed to configure device: %w", err)
	}

	return nil
}
