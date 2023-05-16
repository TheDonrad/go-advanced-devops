// Package encryptionKeyGenerator создаёт RSA ключ
// и записывает публичный и приватный ключ в соот. файлы
// Пример: encryptionKeyGenerator -public-key .\public.key -private-key .\private.key
package main

import (
	"flag"

	"goAdvancedTpl/internal/fabric/encryption"
	"goAdvancedTpl/internal/fabric/logs"
)

func main() {

	var publicKeyPath string
	flag.StringVar(&publicKeyPath, "public-key", "", "public key")
	var privateKeyPath string
	flag.StringVar(&privateKeyPath, "private-key", "", "private key")
	flag.Parse()

	if publicKeyPath == "" || privateKeyPath == "" {
		logs.New().Println("flags \"public-key\" and \"private-key\" must be filled")
		return
	}

	if err := encryption.GenerateKeys(publicKeyPath, privateKeyPath); err != nil {
		logs.New().Println(err.Error())
		return
	}

}
