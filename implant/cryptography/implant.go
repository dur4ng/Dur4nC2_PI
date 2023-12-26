package cryptography

import (
	"errors"
)

var (
	// ECCPublicKey - The implant's ECC public key
	ECCPublicKey = ""
	// eccPrivateKey - The implant's ECC private key
	ECCPrivateKey = ""
	// eccPublicKeySignature - The implant's public key minisigned'd
	ECCPublicKeySignature = ``
	// eccServerPublicKey - Server's ECC public key
	ECCServerPublicKey = ""
	// ErrInvalidPeerKey - Peer to peer key exchange failed
	ErrInvalidPeerKey = errors.New("invalid peer key")
)

func SetSecrets(newEccPublicKey, newEccPrivateKey, newEccPublicKeySignature, newEccServerPublicKey string) {
	ECCPublicKey = newEccPublicKey
	ECCPrivateKey = newEccPrivateKey
	ECCPublicKeySignature = newEccPublicKeySignature
	ECCServerPublicKey = newEccServerPublicKey
}
