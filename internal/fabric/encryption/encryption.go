package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"errors"
	"os"

	"goAdvancedTpl/internal/fabric/logs"
)

func Encrypt(keyPath string, msg []byte) ([]byte, error) {

	key, err := publicKey(keyPath)
	if err != nil {
		return nil, err
	}

	hash := sha512.New()
	return rsa.EncryptOAEP(hash, rand.Reader, key, msg, nil)

}

func publicKey(keyPath string) (*rsa.PublicKey, error) {

	b, err := os.ReadFile(keyPath)
	if err != nil {
		logs.New().Println(err.Error())
		return nil, err
	}

	key, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		logs.New().Println(err.Error())
		return nil, err
	}

	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		err = errors.New("cannot cast type")
		logs.New().Println(err.Error())
		return nil, err
	}

	return pub, nil
}
