package cryptography

import (
	"Dur4nC2/implant/encoders"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
	"sync"
)

var (
	ErrReplayAttack = errors.New("replay attack detected")
	// ErrDecryptFailed
	ErrDecryptFailed = errors.New("decryption failed")
)

// ECCKeyPair - Holds the public/private key pair
type ECCKeyPair struct {
	Public  *[32]byte
	Private *[32]byte
}

// CipherContext - Tracks a series of messages encrypted under the same key
// and detects/prevents replay attacks.
type CipherContext struct {
	Key    [chacha20poly1305.KeySize]byte
	replay *sync.Map
}

// NewCipherContext - Wrapper around creating a cipher context from a key
func NewCipherContext(key [chacha20poly1305.KeySize]byte) *CipherContext {
	return &CipherContext{
		Key:    key,
		replay: &sync.Map{},
	}
}

// RandomKey - Generate random ID of randomIDSize bytes
func RandomKey() [chacha20poly1305.KeySize]byte {
	randBuf := make([]byte, 64)
	rand.Read(randBuf)
	return deriveKeyFrom(randBuf)
}

// deriveKeyFrom - Derives a key from input data using SHA256
func deriveKeyFrom(data []byte) [chacha20poly1305.KeySize]byte {
	digest := sha256.Sum256(data)
	var key [chacha20poly1305.KeySize]byte
	copy(key[:], digest[:chacha20poly1305.KeySize])
	return key
}

// *** Implant - Server cryptography authentication ***
// Encrypt - Encrypt using chacha20poly1305
// https://pkg.go.dev/golang.org/x/crypto/chacha20poly1305
func Encrypt(key [chacha20poly1305.KeySize]byte, plaintext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key[:])
	if err != nil {
		return nil, err
	}
	plaintext = bytes.NewBuffer(encoders.GzipBuf(plaintext)).Bytes()
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(plaintext)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return aead.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt - Decrypt using chacha20poly1305
// https://pkg.go.dev/golang.org/x/crypto/chacha20poly1305
func Decrypt(key [chacha20poly1305.KeySize]byte, ciphertext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key[:])
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aead.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := ciphertext[:aead.NonceSize()], ciphertext[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return encoders.GunzipBuf(plaintext), nil
}

// Decrypt - Decrypt a message with the contextual key and check for replay attacks
func (c *CipherContext) Decrypt(ciphertext []byte) ([]byte, error) {
	plaintext, err := Decrypt(c.Key, ciphertext)
	if err != nil {
		return nil, err
	}
	if 0 < len(ciphertext) {
		digest := sha256.Sum256(ciphertext)
		b64Digest := base64.RawStdEncoding.EncodeToString(digest[:])
		if _, ok := c.replay.LoadOrStore(b64Digest, true); ok {
			return nil, ErrReplayAttack
		}
	}
	return plaintext, nil
}

// Encrypt - Encrypt a message with the contextual key
func (c *CipherContext) Encrypt(plaintext []byte) ([]byte, error) {
	ciphertext, err := Encrypt(c.Key, plaintext)
	if err != nil {
		return nil, err
	}
	if 0 < len(ciphertext) {
		digest := sha256.Sum256(ciphertext)
		b64Digest := base64.RawStdEncoding.EncodeToString(digest[:])
		c.replay.Store(b64Digest, true)
	}
	return ciphertext, nil
}

// *** Implant envelope encryption ***
func (c *CipherContext) EncryptAES(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.Key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}
func (c *CipherContext) DecryptAES(ct []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.Key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := ct[:nonceSize], ct[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
