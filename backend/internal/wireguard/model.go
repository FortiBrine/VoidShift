package wireguard

import "time"

type Network struct {
	ID uint `gorm:"primaryKey"`

	Name       string `gorm:"not null"`
	Address    string `gorm:"not null"`
	ListenPort int    `gorm:"not null"`

	PublicKey  string `gorm:"uniqueIndex;not null"`
	PrivateKey string `gorm:"not null"`

	Peers []Peer `gorm:"foreignKey:NetworkID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Peer struct {
	ID        uint `gorm:"primaryKey"`
	NetworkID uint `gorm:"not null;index"`

	PublicKey  string `gorm:"uniqueIndex;not null"`
	PrivateKey string `gorm:"not null"`

	PresharedKey string
	AllowedIPs   []string `gorm:"serializer:json"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
