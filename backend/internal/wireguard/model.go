package wireguard

type WireguardInterface struct {
	ID string `gorm:"primaryKey"`

	Name       string
	Address    string
	ListenPort *int

	PublicKey  string
	PrivateKey string
}

type Peer struct {
	ID string `gorm:"primaryKey"`

	PublicKey    string
	PresharedKey string
	AllowedIPs   []string `gorm:"type:json"`
}
