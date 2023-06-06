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

type DecrypterServiceTestSuite struct {
	suite.Suite
	privateKeyCertPath string
	publicKeyCertPath  string
	msg                string
	encryptedMsg       string
	privateKey         string
	dec                Decrypter
}

func (suite *DecrypterServiceTestSuite) SetupSuite() {
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

	decrypt, err := NewDecrypter(suite.privateKeyCertPath)

	if err != nil {
		suite.Errorf(err, "decryper create error")
	}

	suite.dec = decrypt

	privateKeyPEM, err := os.ReadFile(suite.privateKeyCertPath)
	if err != nil {
		panic(err)
	}

	suite.privateKey = string(privateKeyPEM)
}

func (suite *DecrypterServiceTestSuite) TearDownSuite() {
	suite.deleteKeys()
	suite.encryptedMsg = ""
	suite.privateKey = ""
	suite.dec = nil
}

func (suite *DecrypterServiceTestSuite) createKeys() {
	manager := NewCertsManager(suite.privateKeyCertPath, suite.publicKeyCertPath)
	if _, _, err := manager.GenerateKeyFiles(); err != nil {
		suite.Errorf(err, "write key file error")
	}
}

func (suite *DecrypterServiceTestSuite) deleteKeys() {
	err := os.Remove(suite.publicKeyCertPath)

	if err != nil {
		suite.Errorf(err, "error removing public file")
	}

	err = os.Remove(suite.privateKeyCertPath)

	if err != nil {
		suite.Errorf(err, "error removing file")
	}
}

func (suite *DecrypterServiceTestSuite) TestGetPrivateKey() {
	key := suite.dec.GetPrivateKey()
	suite.Assert().Equal(suite.privateKey, key)
}

func (suite *DecrypterServiceTestSuite) TestDecryptData() {
	msg, err := suite.dec.DecryptData([]byte(suite.encryptedMsg))

	suite.Assert().NoError(err)
	suite.Assert().Equal(suite.msg, msg)
}

func TestDecrypter(t *testing.T) {
	dec := DecrypterServiceTestSuite{
		privateKeyCertPath: "private.pem",
		publicKeyCertPath:  "public.pem",
		msg:                "Hello World!!!",
	}
	suite.Run(t, &dec)
}
