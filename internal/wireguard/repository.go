package wireguard

import (
	"context"
	"database/sql"

	"github.com/FortiBrine/VoidShift/internal/store"
)

type Repository interface {
	AddNetwork(ctx context.Context, network *Network) error
	GetNetwork(ctx context.Context, networkID uint) (*Network, error)
	GetNetworkWithPeers(ctx context.Context, networkID uint) (*Network, error)
	GetNetworks(ctx context.Context) ([]Network, error)

	GetPeer(ctx context.Context, peerID uint) (Peer, error)
	AddPeer(ctx context.Context, peer *Peer) error
	RemovePeer(ctx context.Context, peerID uint) error
	RemoveNetwork(ctx context.Context, networkID uint) error
}

type SqlcRepository struct {
	db *sql.DB
	q  *store.Queries
}

func NewSqlcRepository(database *sql.DB) *SqlcRepository {
	return &SqlcRepository{db: database, q: store.New(database)}
}

func (r *SqlcRepository) AddNetwork(ctx context.Context, network *Network) error {
	row, err := r.q.CreateNetwork(ctx, store.CreateNetworkParams{
		Name:       network.Name,
		Address:    network.Address,
		ListenPort: int64(network.ListenPort),
		PublicKey:  network.PublicKey,
		PrivateKey: network.PrivateKey,
	})
	if err != nil {
		return err
	}
	network.ID = uint(row.ID)
	return nil
}

func (r *SqlcRepository) GetNetwork(ctx context.Context, networkID uint) (*Network, error) {
	row, err := r.q.GetNetwork(ctx, int64(networkID))
	if err != nil {
		return nil, err
	}
	return new(networkFromDB(row)), nil
}

func (r *SqlcRepository) GetNetworkWithPeers(ctx context.Context, networkID uint) (*Network, error) {
	row, err := r.q.GetNetwork(ctx, int64(networkID))
	if err != nil {
		return nil, err
	}
	n := networkFromDB(row)

	peerRows, err := r.q.GetPeersByNetworkID(ctx, int64(networkID))
	if err != nil {
		return nil, err
	}

	ipRows, err := r.q.GetAllowedIPsByNetworkID(ctx, int64(networkID))
	if err != nil {
		return nil, err
	}
	ipsByPeer := make(map[int64][]string, len(peerRows))
	for _, row := range ipRows {
		ipsByPeer[row.PeerID] = append(ipsByPeer[row.PeerID], row.Ip)
	}

	n.Peers = make([]Peer, len(peerRows))
	for i, p := range peerRows {
		n.Peers[i] = peerFromDB(p, ipsByPeer[p.ID])
	}

	return &n, nil
}

func (r *SqlcRepository) GetNetworks(ctx context.Context) ([]Network, error) {
	rows, err := r.q.GetNetworks(ctx)
	if err != nil {
		return nil, err
	}
	networks := make([]Network, len(rows))
	for i, row := range rows {
		networks[i] = networkFromDB(row)
	}
	return networks, nil
}

func (r *SqlcRepository) GetPeer(ctx context.Context, peerID uint) (Peer, error) {
	row, err := r.q.GetPeer(ctx, int64(peerID))
	if err != nil {
		return Peer{}, err
	}
	ips, err := r.q.GetPeerAllowedIPs(ctx, int64(peerID))
	if err != nil {
		return Peer{}, err
	}
	return peerFromDB(row, ips), nil
}

func (r *SqlcRepository) AddPeer(ctx context.Context, peer *Peer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := store.New(tx)
	row, err := q.CreatePeer(ctx, store.CreatePeerParams{
		NetworkID:    int64(peer.NetworkID),
		PublicKey:    peer.PublicKey,
		PrivateKey:   peer.PrivateKey,
		PresharedKey: peer.PresharedKey,
	})
	if err != nil {
		return err
	}

	for _, ip := range peer.AllowedIPs {
		if err := q.InsertPeerAllowedIP(ctx, store.InsertPeerAllowedIPParams{
			PeerID: row.ID,
			Ip:     ip,
		}); err != nil {
			return err
		}
	}

	peer.ID = uint(row.ID)
	return tx.Commit()
}

func (r *SqlcRepository) RemovePeer(ctx context.Context, peerID uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := store.New(tx)
	if err := q.DeletePeerAllowedIPs(ctx, int64(peerID)); err != nil {
		return err
	}
	if err := q.DeletePeer(ctx, int64(peerID)); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *SqlcRepository) RemoveNetwork(ctx context.Context, networkID uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := store.New(tx)
	if err := q.DeleteAllowedIPsByNetworkID(ctx, int64(networkID)); err != nil {
		return err
	}
	if err := q.DeletePeersByNetworkID(ctx, int64(networkID)); err != nil {
		return err
	}
	if err := q.DeleteNetwork(ctx, int64(networkID)); err != nil {
		return err
	}
	return tx.Commit()
}

func networkFromDB(row store.Network) Network {
	return Network{
		ID:         uint(row.ID),
		Name:       row.Name,
		Address:    row.Address,
		ListenPort: int(row.ListenPort),
		PublicKey:  row.PublicKey,
		PrivateKey: row.PrivateKey,
	}
}

func peerFromDB(row store.Peer, ips []string) Peer {
	if ips == nil {
		ips = []string{}
	}
	return Peer{
		ID:           uint(row.ID),
		NetworkID:    uint(row.NetworkID),
		PublicKey:    row.PublicKey,
		PrivateKey:   row.PrivateKey,
		PresharedKey: row.PresharedKey,
		AllowedIPs:   ips,
	}
}
