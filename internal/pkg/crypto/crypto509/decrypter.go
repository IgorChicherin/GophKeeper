package crypto509

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type Decrypter interface {
	DecryptData(data []byte) (string, error)
	GetPrivateKey() string
}

func NewDecrypter(privateKeyFileName string) (Decrypter, error) {

	privateKeyPEM, err := os.ReadFile(privateKeyFileName)
	if err != nil {
		return nil, err
	}

	return decrypter{privateKey: string(privateKeyPEM)}, nil
}

type decrypter struct {
	privateKey string
}

func (d decrypter) DecryptData(data []byte) (string, error) {
	privateKeyBlock, _ := pem.Decode([]byte(d.privateKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)

	if err != nil {
		return "", err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func (d decrypter) GetPrivateKey() string {
	return d.privateKey
}
