package crypto509

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EncrypterServiceTestSuite struct {
	suite.Suite
	privateKeyCertPath string
	publicKeyCertPath  string
	msg                string
	encryptedMsg       string
	publicKey          string
	enc                Encrypter
}

func (suite *EncrypterServiceTestSuite) SetupTest() {
	suite.Assert().NotEmpty(suite.privateKeyCertPath)
	suite.Assert().NotEmpty(suite.publicKeyCertPath)
	suite.Assert().NotEmpty(suite.msg)

	suite.createKeys()

	publicKeyPEM, err := os.ReadFile(suite.publicKeyCertPath)
	if err != nil {
		panic(err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)

	k, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	plaintext := []byte("Hello World!!!")
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, k.(*rsa.PublicKey), plaintext)

	if err != nil {
		panic(err)
	}

	suite.encryptedMsg = string(ciphertext)

	encrypt, err := NewEncrypter(suite.publicKeyCertPath)

	if err != nil {
		suite.Errorf(err, "encrypter create error")
	}

	suite.enc = encrypt

	suite.publicKey = string(publicKeyPEM)
}

func (suite *EncrypterServiceTestSuite) TearDownSuite() {
	suite.deleteKeys()
	suite.encryptedMsg = ""
	suite.publicKey = ""
	suite.enc = nil
}

func (suite *EncrypterServiceTestSuite) createKeys() {
	manager := NewCertsManager(suite.privateKeyCertPath, suite.publicKeyCertPath)
	if _, _, err := manager.GenerateKeyFiles(); err != nil {
		suite.Errorf(err, "write key file error")
	}
}

func (suite *EncrypterServiceTestSuite) deleteKeys() {
	err := os.Remove(suite.publicKeyCertPath)

	if err != nil {
		suite.Errorf(err, "error removing public file")
	}

	err = os.Remove(suite.privateKeyCertPath)

	if err != nil {
		suite.Errorf(err, "error removing file")
	}
}

func (suite *EncrypterServiceTestSuite) TestGetPublicKey() {
	key := suite.enc.GetPublicKey()

	suite.Assert().NotEmpty(key)
	suite.Assert().Equal(suite.publicKey, key)
}

func (suite *EncrypterServiceTestSuite) TestEncryptData() {
	encData, err := suite.enc.EncryptData(suite.msg)
	suite.Assert().Equal(err, nil)

	newDec, err := NewDecrypter(suite.privateKeyCertPath)
	suite.Assert().Equal(err, nil)

	decMsg, err := newDec.DecryptData([]byte(encData))
	suite.Assert().Equal(err, nil)

	decMsgFromEnc, err := newDec.DecryptData([]byte(suite.encryptedMsg))
	suite.Assert().Equal(err, nil)

	suite.Assert().Equal(suite.msg, decMsg)
	suite.Assert().Equal(suite.msg, decMsgFromEnc)
	suite.Assert().Equal(decMsg, decMsgFromEnc)
}

func TestEncrypter(t *testing.T) {
	enc := EncrypterServiceTestSuite{
		privateKeyCertPath: "private.pem",
		publicKeyCertPath:  "public.pem",
		msg:                "Hello World!!!",
	}
	suite.Run(t, &enc)
}
