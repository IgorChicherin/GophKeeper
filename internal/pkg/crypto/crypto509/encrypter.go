package crypto509

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type Encrypter interface {
	EncryptData(data string) (string, error)
	GetPublicKey() string
}

func NewEncrypter(publicKeyFileName string) (Encrypter, error) {
	publicKeyPEM, err := os.ReadFile(publicKeyFileName)
	if err != nil {
		return nil, err
	}

	return encrypter{publicKey: string(publicKeyPEM)}, nil
}

type encrypter struct {
	publicKey string
}

func (e encrypter) EncryptData(data string) (string, error) {
	publicKeyBlock, _ := pem.Decode([]byte(e.publicKey))

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return "", err
	}

	plaintext := []byte(data)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), plaintext)

	if err != nil {
		return "", err
	}

	return string(ciphertext), nil
}

func (e encrypter) GetPublicKey() string {
	return e.publicKey
}
