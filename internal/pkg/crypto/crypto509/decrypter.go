package crypto509

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type Decrypter interface {
	DecryptData(data []byte) ([]byte, error)
	GetPrivateKey() string
}

func NewDecrypter(privateKey []byte) (Decrypter, error) {
	return decrypter{privateKey: string(privateKey)}, nil
}

type decrypter struct {
	privateKey string
}

func (d decrypter) DecryptData(data []byte) ([]byte, error) {
	privateKeyBlock, _ := pem.Decode([]byte(d.privateKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)

	if err != nil {
		return []byte{}, err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return []byte{}, err
	}
	return plaintext, nil
}

func (d decrypter) GetPrivateKey() string {
	return d.privateKey
}
