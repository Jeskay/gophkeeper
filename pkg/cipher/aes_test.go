package cipher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesCipher(t *testing.T) {
	data := []byte("Know thy self, know thy enemy, and in a hundred battles you will never be in peril")
	key := []byte("N1PCdw3M2B1TfJhoaY2mL736p2vCUc47")
	aesCipher, err := NewAESCipher(key)
	assert.NoErrorf(t, err, "failed to init")
	encrypted, err := aesCipher.Encrypt(data)
	assert.NoErrorf(t, err, "error during encryption")
	decrypted, err := aesCipher.Decrypt(encrypted)
	if assert.NoErrorf(t, err, "error during decryption") {
		assert.Equal(t, data, decrypted)
		assert.NotEqual(t, data, encrypted)
	}
}
