package wireguard

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/FortiBrine/VoidShift/internal/shared"
	"github.com/skip2/go-qrcode"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gorm.io/gorm"
)

type Service struct {
	repository  Repository
	client      *wgctrl.Client
	hostAddress string
}

func NewService(repository Repository, client *wgctrl.Client, hostAddress string) *Service {
	return &Service{
		repository:  repository,
		client:      client,
		hostAddress: hostAddress,
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
	allowedIPs []string,
) (*Peer, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	publicKey := privateKey.PublicKey()
	psk, err := wgtypes.GenerateKey()

	if err != nil {
		return nil, fmt.Errorf("failed to generate preshared key: %w", err)
	}

	peer := &Peer{
		NetworkID: networkID,

		PublicKey:  publicKey.String(),
		PrivateKey: privateKey.String(),

		PresharedKey: psk.String(),
		AllowedIPs:   allowedIPs,
	}

	return peer, s.repository.AddPeer(ctx, peer)
}

func (s *Service) RemovePeer(
	ctx context.Context,
	peerID uint,
) error {
	peer, err := s.repository.GetPeer(ctx, peerID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared.ErrPeerNotFound
		}

		return err
	}

	network, err := s.repository.GetNetwork(ctx, peer.NetworkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared.ErrPeerNotFound
		}

		return err
	}

	up, err := IsDeviceUp(network.Name)
	if err != nil {
		return fmt.Errorf("failed to check device state: %w", err)
	}

	if !up {
		_, err = s.repository.RemovePeer(ctx, peerID)

		return err
	}

	publicKey, err := wgtypes.ParseKey(peer.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to parse peer public key: %w", err)
	}

	err = s.client.ConfigureDevice(network.Name, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey: publicKey,
				Remove:    true,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to remove peer from device: %w", err)
	}

	_, err = s.repository.RemovePeer(ctx, peerID)

	return err
}

func (s *Service) GetPeerConfig(
	ctx context.Context,
	peerID uint,
) (string, error) {
	peer, err := s.repository.GetPeer(ctx, peerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", shared.ErrPeerNotFound
		}

		return "", err
	}

	network, err := s.repository.GetNetwork(ctx, peer.NetworkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", shared.ErrNetworkNotFound
		}

		return "", err
	}

	_, mask, found := strings.Cut(network.Address, "/")
	if !found {
		return "", shared.ErrNetworkNotFound
	}

	config := strings.Builder{}
	config.WriteString("[Interface]\n")
	config.WriteString(fmt.Sprintf("PrivateKey = %s\n", peer.PrivateKey))
	config.WriteString(fmt.Sprintf("DNS = %s\n", "1.1.1.1"))

	var processed []string

	for _, ip := range peer.AllowedIPs {
		if !strings.Contains(ip, "/") {
			ip = ip + "/" + mask
		}
		processed = append(processed, ip)
	}

	addresses := strings.Join(processed, ", ")
	if addresses != "" {
		config.WriteString(fmt.Sprintf("Address = %s\n", addresses))
	}

	config.WriteString("\n[Peer]\n")
	config.WriteString(fmt.Sprintf("PublicKey = %s\n", network.PublicKey))

	if peer.PresharedKey != "" {
		config.WriteString(fmt.Sprintf("PresharedKey = %s\n", peer.PresharedKey))
	}

	endpoint := net.JoinHostPort(s.hostAddress, strconv.Itoa(network.ListenPort))
	config.WriteString(fmt.Sprintf("Endpoint = %s\n", endpoint))
	config.WriteString(fmt.Sprintf("AllowedIPs = %s, %s\n", "0.0.0.0/0", "::/0"))

	return config.String(), nil
}

func (s *Service) GetPeerConfigQR(
	ctx context.Context,
	peerID uint,
) ([]byte, error) {
	config, err := s.GetPeerConfig(ctx, peerID)
	if err != nil {
		return nil, err
	}

	qrCode, err := qrcode.Encode(config, qrcode.Medium, 512)
	if err != nil {
		return nil, fmt.Errorf("failed to generate qr code: %w", err)
	}

	return qrCode, nil
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
) (*Network, error) {
	return s.repository.GetNetwork(ctx, networkID)
}

func (s *Service) GetNetworkWithPeers(
	ctx context.Context,
	networkID uint,
) (*Network, error) {
	network, err := s.repository.GetNetworkWithPeers(ctx, networkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNetworkNotFound
		}

		return nil, err
	}

	return network, err
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
			_, ipNet, err := net.ParseCIDR(allowedIP + "/32")
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

func (s *Service) DownNetwork(
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

	return nil
}
