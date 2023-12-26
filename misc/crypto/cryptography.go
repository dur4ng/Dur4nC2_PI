package cryptography

import (
	"Dur4nC2/misc/compress"
	"Dur4nC2/server/db"
	_serverRepository "Dur4nC2/server/domain/server/repository/postgres"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/nacl/box"
	"gorm.io/gorm"
	"io"
	"math/big"
	"sync"
)

const (
	ServerECCKeyPairKey = "server.ecc"
)

var (
	// ErrInvalidKeyLength - Invalid key length
	ErrInvalidKeyLength = errors.New("invalid length")

	// ErrReplayAttack - Replay attack
	ErrReplayAttack = errors.New("replay attack detected")

	// ErrDecryptFailed
	ErrDecryptFailed = errors.New("decryption failed")
)

var (
	serverRepository = _serverRepository.NewPostgresServerRepository(db.Session())
)

// deriveKeyFrom - Derives a key from input data using SHA256
func deriveKeyFrom(data []byte) [chacha20poly1305.KeySize]byte {
	digest := sha256.Sum256(data)
	var key [chacha20poly1305.KeySize]byte
	copy(key[:], digest[:chacha20poly1305.KeySize])
	return key
}

// RandomKey - Generate random ID of randomIDSize bytes
func RandomKey() [chacha20poly1305.KeySize]byte {
	randBuf := make([]byte, 64)
	rand.Read(randBuf)
	return deriveKeyFrom(randBuf)
}

// KeyFromBytes - Convert to fixed length buffer
func KeyFromBytes(data []byte) ([chacha20poly1305.KeySize]byte, error) {
	var key [chacha20poly1305.KeySize]byte
	if len(data) != chacha20poly1305.KeySize {
		// We cannot return nil due to the fixed length buffer type ...
		// and it seems like a really bad idea to return a zero key in case
		// the error is not checked by the caller, so instead we return a
		// random key, which should break everything if the error is not checked.
		return RandomKey(), ErrInvalidKeyLength
	}
	copy(key[:], data)
	return key, nil
}

// ECCKeyPair - Holds the public/private key pair
type ECCKeyPair struct {
	Public  *[32]byte `json:"public"`
	Private *[32]byte `json:"private"`
}

// PublicBase64 - Base64 encoded public key
func (e *ECCKeyPair) PublicBase64() string {
	return base64.RawStdEncoding.EncodeToString(e.Public[:])
}

// PrivateBase64 - Base64 encoded private key
func (e *ECCKeyPair) PrivateBase64() string {
	return base64.RawStdEncoding.EncodeToString(e.Private[:])
}

// RandomECCKeyPair - Generate a random Curve 25519 key pair
func RandomECCKeyPair() (*ECCKeyPair, error) {
	public, private, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &ECCKeyPair{Public: public, Private: private}, nil
}

// ECCServerKeyPair - Get the server's ECC key pair
func ECCServerKeyPair() *ECCKeyPair {
	data, err := serverRepository.GetKeyValue(ServerECCKeyPairKey)
	if err == gorm.ErrRecordNotFound {
		keyPair, err := generateServerECCKeyPair()
		if err != nil {
			panic(err)
		}
		return keyPair
	}
	keyPair := &ECCKeyPair{}
	err = json.Unmarshal([]byte(data), keyPair)
	if err != nil {
		panic(err)
	}
	return keyPair

}
func generateServerECCKeyPair() (*ECCKeyPair, error) {
	keyPair, err := RandomECCKeyPair()
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(keyPair)
	if err != nil {
		return nil, err
	}
	err = serverRepository.SetKeyValue(ServerECCKeyPairKey, string(data))
	return keyPair, err
}

// *** Implant - Server cryptography authentication ***
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

// Encrypt - Encrypt using chacha20poly1305
// https://pkg.go.dev/golang.org/x/crypto/chacha20poly1305
func Encrypt(key [chacha20poly1305.KeySize]byte, plaintext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key[:])
	if err != nil {
		return nil, err
	}
	plaintext = bytes.NewBuffer(compressors.GzipBuf(plaintext)).Bytes()
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
	return compressors.GunzipBuf(plaintext), nil
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

// *** stage2 crypto ***
// PreludeEncrypt the results
func PreludeEncrypt(data []byte, key []byte, iv []byte) []byte {
	plainText, err := pad(data, aes.BlockSize)
	if err != nil {
		return make([]byte, 0)
	}
	block, _ := aes.NewCipher(key)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	// Create a random IV if none was provided
	// len(nil) returns 0
	if len(iv) == 0 {
		iv = cipherText[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return make([]byte, 0)
		}
	} else {
		// make sure we copy the IV
		copy(cipherText[:aes.BlockSize], iv)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)
	return cipherText
}

// PreludeDecrypt a command
func PreludeDecrypt(data []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	data, _ = unpad(data, aes.BlockSize)
	return data
}
func pad(buf []byte, size int) ([]byte, error) {
	bufLen := len(buf)
	padLen := size - bufLen%size
	padded := make([]byte, bufLen+padLen)
	copy(padded, buf)
	for i := 0; i < padLen; i++ {
		padded[bufLen+i] = byte(padLen)
	}
	return padded, nil
}
func unpad(padded []byte, size int) ([]byte, error) {
	if len(padded)%size != 0 {
		return nil, errors.New("pkcs7: Padded value wasn't in correct size")
	}
	bufLen := len(padded) - int(padded[len(padded)-1])
	buf := make([]byte, bufLen)
	copy(buf, padded[:bufLen])
	return buf, nil
}

func GenerateRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	max := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, max)
		result[i] = chars[n.Int64()]
	}

	return string(result)
}
