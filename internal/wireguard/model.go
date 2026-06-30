package wireguard

type Network struct {
	ID         uint
	Name       string
	Address    string
	ListenPort int
	PublicKey  string
	PrivateKey string
	Peers      []Peer
}

type Peer struct {
	ID           uint
	NetworkID    uint
	PublicKey    string
	PrivateKey   string
	PresharedKey string
	AllowedIPs   []string
}
