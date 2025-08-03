package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type aesCipher struct {
	secretKey []byte
	gcm       cipher.AEAD
}

func NewAESCipher(secretKey []byte) (*aesCipher, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &aesCipher{secretKey: secretKey, gcm: gcm}, nil
}

func (w *aesCipher) Encrypt(data []byte) ([]byte, error) {

	nonce := make([]byte, w.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return w.gcm.Seal(nonce, nonce, data, nil), nil //TODO: add salt data
}

func (w *aesCipher) Decrypt(data []byte) ([]byte, error) {
	nonce := data[:w.gcm.NonceSize()]
	ciphered := data[w.gcm.NonceSize():]
	return w.gcm.Open(nil, nonce, ciphered, nil)
}
