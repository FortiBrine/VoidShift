package wireguard

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/skip2/go-qrcode"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	ErrNetworkNotFound = errors.New("network not found")
	ErrPeerNotFound    = errors.New("peer not found")
)

type Service struct {
	repository  Repository
	client      *wgctrl.Client
	hostAddress string
	logger      *slog.Logger

	mu sync.Mutex
}

func NewService(
	repository Repository,
	hostAddress string,
	logger *slog.Logger,
) (*Service, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	return &Service{
		repository:  repository,
		client:      client,
		hostAddress: hostAddress,
		logger:      logger,
	}, nil
}

func (s *Service) Load(ctx context.Context) error {
	networks, err := s.repository.GetEnabledNetworks(ctx)
	if err != nil {
		return fmt.Errorf("getting enabled networks: %w", err)
	}

	for _, network := range networks {
		if err := s.UpNetwork(ctx, network.ID); err != nil {
			s.logger.Error("restoring network on startup",
				slog.String("network", network.Name),
				slog.Any("error", err),
			)
		}
	}

	return nil
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
		return nil, fmt.Errorf("generating private key: %w", err)
	}

	publicKey := privateKey.PublicKey()
	psk, err := wgtypes.GenerateKey()

	if err != nil {
		return nil, fmt.Errorf("generating preshared key: %w", err)
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
		return err
	}

	network, err := s.repository.GetNetwork(ctx, peer.NetworkID)
	if err != nil {
		if errors.Is(err, ErrNetworkNotFound) {
			return ErrPeerNotFound
		}

		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	iface := IfaceName(network.Name, network.ID)
	up, err := IsDeviceUp(iface)
	if err != nil {
		return fmt.Errorf("checking device state: %w", err)
	}

	if !up {
		return s.repository.RemovePeer(ctx, peerID)
	}

	publicKey, err := wgtypes.ParseKey(peer.PublicKey)
	if err != nil {
		return fmt.Errorf("parsing peer public key: %w", err)
	}

	err = s.client.ConfigureDevice(iface, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey: publicKey,
				Remove:    true,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("removing peer from device: %w", err)
	}

	return s.repository.RemovePeer(ctx, peerID)
}

func (s *Service) GetPeerConfig(
	ctx context.Context,
	peerID uint,
) (string, error) {
	peer, err := s.repository.GetPeer(ctx, peerID)
	if err != nil {
		return "", err
	}

	network, err := s.repository.GetNetwork(ctx, peer.NetworkID)
	if err != nil {
		return "", err
	}

	_, mask, found := strings.Cut(network.Address, "/")
	if !found {
		return "", ErrNetworkNotFound
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
		return nil, fmt.Errorf("generating qr code: %w", err)
	}

	return qrCode, nil
}

func (s *Service) RemoveNetwork(
	ctx context.Context,
	networkID uint,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	network, err := s.repository.GetNetwork(ctx, networkID)
	if err != nil {
		return err
	}
	iface := IfaceName(network.Name, network.ID)
	if err = RemoveDevice(iface); err != nil {
		return fmt.Errorf("removing device: %w", err)
	}

	return s.repository.RemoveNetwork(ctx, networkID)
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
	return s.repository.GetNetworkWithPeers(ctx, networkID)
}

func (s *Service) UpNetwork(
	ctx context.Context,
	networkID uint,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	network, err := s.repository.GetNetworkWithPeers(ctx, networkID)
	if err != nil {
		return err
	}

	iface := IfaceName(network.Name, network.ID)
	if err := CreateDevice(iface); err != nil {
		return fmt.Errorf("creating device: %w", err)
	}

	if err := SetDeviceAddress(iface, network.Address); err != nil {
		return fmt.Errorf("setting device address: %w", err)
	}

	privateKey, err := wgtypes.ParseKey(network.PrivateKey)
	if err != nil {
		return fmt.Errorf("parsing private network key: %w", err)
	}

	peers := make([]wgtypes.PeerConfig, len(network.Peers))
	for i, peer := range network.Peers {
		publicKey, err := wgtypes.ParseKey(peer.PublicKey)
		if err != nil {
			return fmt.Errorf("parsing public peer key: %w", err)
		}

		var presharedKey *wgtypes.Key
		if peer.PresharedKey != "" {
			psk, err := wgtypes.ParseKey(peer.PresharedKey)
			if err != nil {
				return fmt.Errorf("parsing preshared key: %w", err)
			}
			presharedKey = &psk
		}

		allowedIPs := make([]net.IPNet, len(peer.AllowedIPs))

		for j, allowedIP := range peer.AllowedIPs {
			_, ipNet, err := net.ParseCIDR(allowedIP + "/32")
			if err != nil {
				return fmt.Errorf("parsing allowed IP %q: %w", allowedIP, err)
			}

			allowedIPs[j] = *ipNet
		}

		peers[i] = wgtypes.PeerConfig{
			PublicKey:    publicKey,
			PresharedKey: presharedKey,
			AllowedIPs:   allowedIPs,
		}
	}

	if err := s.client.ConfigureDevice(iface, wgtypes.Config{
		PrivateKey:   &privateKey,
		ListenPort:   &network.ListenPort,
		ReplacePeers: true,
		Peers:        peers,
	}); err != nil {
		return fmt.Errorf("configuring device: %w", err)
	}

	return s.repository.SetNetworkEnabled(ctx, networkID, true)
}

func (s *Service) DownNetwork(
	ctx context.Context,
	networkID uint,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	network, err := s.repository.GetNetwork(ctx, networkID)
	if err != nil {
		return err
	}

	iface := IfaceName(network.Name, network.ID)
	if err = RemoveDevice(iface); err != nil {
		return fmt.Errorf("removing device: %w", err)
	}

	return s.repository.SetNetworkEnabled(ctx, networkID, false)
}

func (s *Service) GetStats(
	ctx context.Context,
) (Stats, error) {
	networkCount, err := s.repository.CountNetworks(ctx)
	if err != nil {
		return Stats{}, fmt.Errorf("counting networks: %w", err)
	}

	peerCount, err := s.repository.CountPeers(ctx)
	if err != nil {
		return Stats{}, fmt.Errorf("counting peers: %w", err)
	}

	return Stats{NetworkCount: networkCount, PeerCount: peerCount}, nil
}

func (s *Service) Close() error {
	return s.client.Close()
}
