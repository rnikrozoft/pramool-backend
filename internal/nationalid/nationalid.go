package nationalid

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strings"
)

var (
	ErrInvalidKey = errors.New("national id encryption key must be 32 bytes base64")
	ErrNotReady   = errors.New("national id codec not initialized")
)

var codec *Codec

type Codec struct {
	gcm cipher.AEAD
}

func InitFromBase64Key(keyB64 string) error {
	keyB64 = strings.TrimSpace(keyB64)
	if keyB64 == "" {
		return ErrInvalidKey
	}
	raw, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil || len(raw) != 32 {
		return ErrInvalidKey
	}
	block, err := aes.NewCipher(raw)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	codec = &Codec{gcm: gcm}
	return nil
}

func Hash(plain string) string {
	plain = strings.TrimSpace(plain)
	if plain == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

func Encrypt(plain string) (string, error) {
	if codec == nil {
		return "", ErrNotReady
	}
	plain = strings.TrimSpace(plain)
	if plain == "" {
		return "", nil
	}
	nonce := make([]byte, codec.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := codec.gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(enc string) (string, error) {
	if codec == nil {
		return "", ErrNotReady
	}
	enc = strings.TrimSpace(enc)
	if enc == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	nonceSize := codec.gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}
	plain, err := codec.gcm.Open(nil, raw[:nonceSize], raw[nonceSize:], nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func PrepareStorage(plain string) (hash string, enc string, err error) {
	plain = strings.TrimSpace(plain)
	if plain == "" {
		return "", "", nil
	}
	hash = Hash(plain)
	enc, err = Encrypt(plain)
	return hash, enc, err
}
