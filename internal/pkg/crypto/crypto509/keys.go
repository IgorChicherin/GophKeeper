package crypto509

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

type CertsManager interface {
	GetCerts() ([]byte, []byte, error)
	LoadKeysFromFiles() ([]byte, []byte, error)
	GenerateKeyFiles() ([]byte, []byte, error)
	CreateKeyFile(keyFileName string, keyPEM []byte) error
	CreateKeysPEM() ([]byte, []byte, error)
}

type certsManager struct {
	PrivateKeyCertPath string
	PublicKeyCertPath  string
}

func NewCertsManager(privateKeyCertPath, publicKeyCertPath string) CertsManager {
	return certsManager{PrivateKeyCertPath: privateKeyCertPath, PublicKeyCertPath: publicKeyCertPath}
}

func (m certsManager) GetCerts() ([]byte, []byte, error) {
	if !fileExists(m.PrivateKeyCertPath) && !fileExists(m.PrivateKeyCertPath) {
		private, public, err := m.GenerateKeyFiles()
		if err != nil {
			log.Errorln(err)
			return nil, nil, err
		}
		return private, public, err
	}

	return m.LoadKeysFromFiles()
}

func (m certsManager) LoadKeysFromFiles() ([]byte, []byte, error) {
	publicKeyPEM, err := os.ReadFile(m.PublicKeyCertPath)
	if err != nil {
		log.Errorln(err)
		return nil, nil, err
	}

	privateKeyPEM, err := os.ReadFile(m.PrivateKeyCertPath)
	if err != nil {
		log.Errorln(err)
		return nil, nil, err
	}
	return privateKeyPEM, publicKeyPEM, nil
}

func (m certsManager) GenerateKeyFiles() ([]byte, []byte, error) {
	privateKeyPEM, publicKeyPEM, err := m.CreateKeysPEM()

	err = m.CreateKeyFile(m.PrivateKeyCertPath, privateKeyPEM)
	if err != nil {
		log.Errorln(err)
		return []byte{}, []byte{}, err
	}

	err = m.CreateKeyFile(m.PublicKeyCertPath, publicKeyPEM)
	if err != nil {
		log.Errorln(err)
		return []byte{}, []byte{}, err
	}
	return privateKeyPEM, publicKeyPEM, err
}

func (m certsManager) CreateKeyFile(keyFileName string, keyPEM []byte) error {
	if err := os.WriteFile(keyFileName, keyPEM, 0644); err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (m certsManager) CreateKeysPEM() ([]byte, []byte, error) {
	keyPair, err := generateKeyPair()
	if err != nil {
		log.Errorln(err)
		return nil, nil, err
	}

	publicKeyPEM, err := getPublicKeyPEM(keyPair)
	if err != nil {
		log.Errorln(err)
		return nil, nil, err
	}

	privateKeyPEM := getPrivateKeyPEM(keyPair)

	return privateKeyPEM, publicKeyPEM, nil
}

func generateKeyPair() (rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Errorln(err, "generate private key error")
		return rsa.PrivateKey{}, err
	}

	return *privateKey, nil
}

func getPrivateKeyPEM(keyPair rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(&keyPair)
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
}

func getPublicKeyPEM(keyPair rsa.PrivateKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&keyPair.PublicKey)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}), nil
}

func fileExists(filePath string) bool {
	if info, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) || info.IsDir() {
		return false
	}
	return true
}
