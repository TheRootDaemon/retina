package model

import (
	"gorm.io/gorm"
)

// Host represents a remote SSH host managed by hop.
type Host struct {
	gorm.Model

	// Alias is a unique, user-friendly name for the host.
	Alias string `gorm:"uniqueIndex;not null"`

	// Hostname is the IP address or DNS name of the remote host.
	Hostname string `gorm:"not null"`

	// Port is the SSH port.
	Port uint16 `gorm:"not null;default:22"`

	// Username is the SSH login name.
	Username string `gorm:"not null"`

	// Password contains the encrypted SSH password.
	Password []byte

	// PrivateKey contains the encrypted private key.
	PrivateKey []byte

	// PublicKey contains the SSH public key.
	PublicKey []byte
}
