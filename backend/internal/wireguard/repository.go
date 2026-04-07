package wireguard

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Migrate() error
	AddNetwork(ctx context.Context, network *Network) error

	GetNetwork(ctx context.Context, networkID uint) (Network, error)
	GetNetworkWithPeers(ctx context.Context, networkID uint) (Network, error)
	GetNetworks(ctx context.Context) ([]Network, error)

	AddPeer(ctx context.Context, peer *Peer) error
	RemovePeer(ctx context.Context, peerID uint) (int, error)

	UpdateNetwork(ctx context.Context, networkID uint, network Network) (int, error)
	RemoveNetwork(ctx context.Context, networkID uint) (int, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Migrate() error {
	return r.db.AutoMigrate(&Network{}, &Peer{})
}

func (r *GormRepository) AddNetwork(ctx context.Context, network *Network) error {
	return gorm.G[Network](r.db).Create(ctx, network)
}

func (r *GormRepository) GetNetwork(ctx context.Context, networkID uint) (Network, error) {
	return gorm.G[Network](r.db).
		Where("id = ?", networkID).
		First(ctx)
}

func (r *GormRepository) GetNetworkWithPeers(ctx context.Context, networkID uint) (Network, error) {
	return gorm.G[Network](r.db).
		Joins(clause.JoinTarget{
			Association: "Peers",
		}, nil).
		Where("id = ?", networkID).First(ctx)
}

func (r *GormRepository) GetNetworks(ctx context.Context) ([]Network, error) {
	return gorm.G[Network](r.db).Find(ctx)
}

func (r *GormRepository) AddPeer(ctx context.Context, peer *Peer) error {
	return gorm.G[Peer](r.db).
		Create(ctx, peer)
}

func (r *GormRepository) RemovePeer(ctx context.Context, peerID uint) (int, error) {
	return gorm.G[Peer](r.db).
		Where("id = ?", peerID).
		Delete(ctx)
}

func (r *GormRepository) UpdateNetwork(
	ctx context.Context,
	networkID uint,
	network Network,
) (int, error) {
	return gorm.G[Network](r.db).
		Where("id = ?", networkID).
		Updates(ctx, network)
}

func (r *GormRepository) RemoveNetwork(ctx context.Context, networkID uint) (int, error) {
	return gorm.G[Network](r.db).
		Where("id = ?", networkID).
		Delete(ctx)
}
