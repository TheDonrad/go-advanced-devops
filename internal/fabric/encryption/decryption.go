package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"os"

	"goAdvancedTpl/internal/fabric/logs"
)

func Decrypt(keyPath string, msg []byte) ([]byte, error) {

	key, err := privateKey(keyPath)
	if err != nil {
		return nil, err
	}

	hash := sha512.New()
	return rsa.DecryptOAEP(hash, rand.Reader, key, msg, nil)

}

func privateKey(keyPath string) (*rsa.PrivateKey, error) {

	b, err := os.ReadFile(keyPath)
	if err != nil {
		logs.Logger().Println(err.Error())
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		logs.Logger().Println(err.Error())
		return nil, err
	}

	return key, nil

}
