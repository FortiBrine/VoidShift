package wireguard

import (
	"context"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
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

	return network, s.repository.AddWireGuardNetwork(ctx, network)
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
	_, err := s.repository.RemoveNetwork(ctx, networkID)
	return err
}
