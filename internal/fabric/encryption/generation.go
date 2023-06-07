package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"goAdvancedTpl/internal/fabric/logs"
)

func GenerateKeys(publicKeyPath string, privateKeyPath string) error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logs.Logger().Println(err.Error())
	}

	err = writePublicKey(publicKeyPath, &privateKey.PublicKey)
	if err != nil {
		return err
	}

	err = writePrivateKey(privateKeyPath, privateKey)
	if err != nil {
		return err
	}

	return nil

}

func writePublicKey(filePath string, key *rsa.PublicKey) error {
	pubPKI, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		logs.Logger().Println(err.Error())
		return err
	}

	pubEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubPKI,
	})

	if err = os.WriteFile(filePath, pubEncoded, 0666); err != nil {
		logs.Logger().Println(err.Error())
		return err
	}

	return nil
}

func writePrivateKey(filePath string, key *rsa.PrivateKey) error {
	privateEncoded := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	if err := os.WriteFile(filePath, privateEncoded, 0666); err != nil {
		logs.Logger().Println(err.Error())
		return err
	}

	return nil
}
