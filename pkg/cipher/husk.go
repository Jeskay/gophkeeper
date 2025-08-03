package cipher

type huskCipher struct{}

func NewHuskCipher() *huskCipher {
	return &huskCipher{}
}

func (*huskCipher) Encrypt(data []byte) ([]byte, error) {
	return data, nil
}

func (*huskCipher) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}
