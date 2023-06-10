package crypto509

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type Encrypter interface {
	EncryptData(data []byte) (string, error)
	GetPublicKey() string
}

func NewEncrypter(publicKey []byte) (Encrypter, error) {
	return encrypter{publicKey: string(publicKey)}, nil
}

type encrypter struct {
	publicKey string
}

func (e encrypter) EncryptData(data []byte) (string, error) {
	publicKeyBlock, _ := pem.Decode([]byte(e.publicKey))

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return "", err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), data)

	if err != nil {
		return "", err
	}

	return string(ciphertext), nil
}

func (e encrypter) GetPublicKey() string {
	return e.publicKey
}
