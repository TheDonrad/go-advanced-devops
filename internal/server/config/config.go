// Package config служит для определения настроек серевера сбора метрик
package config

import (
	"flag"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

// SettingsList хранит настройки сервера
// Если заполнен параметр DBConnString, считаем, что в качестве хранилища следует использовать БД
// Иначе - файл
type SettingsList struct {
	Addr          string        // Адрес для получения метрик
	Key           string        // Ключ для отправки шифрованного хеша метрики по алгоритму sha256
	DBConnString  string        // Строка соединения с БД
	StoreFile     string        // Путь к файлу для хранения метрик
	StoreInterval time.Duration // Период сохранения настроек
	Restore       bool          // Восстанавливать метрики из хранилища при запуске
	CryptoKey     string        // Ключ шифрования

}

// Config возвращает настройки агента из переменных окружения или флагов запуска.
// У переменных окружения приоритет перед флагами
func Config() *SettingsList {
	settings := SettingsList{
		Addr:          "127.0.0.1:8080",
		StoreInterval: 120 * time.Second,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       true,
		DBConnString:  "",
		CryptoKey:     "",
	}
	settings.setConfigFlags()
	settings.setConfigEnv()

	return &settings
}

func (settings *SettingsList) setConfigFlags() {

	flag.StringVar(&settings.Addr, "a", settings.Addr, "host to listen on")

	flag.StringVar(&settings.StoreFile, "f", settings.StoreFile, "file to store metrics")

	flag.Func("i", "store interval", func(flagValue string) error {
		settings.StoreInterval, _ = time.ParseDuration(flagValue)
		return nil
	})
	flag.BoolVar(&settings.Restore, "r", settings.Restore, "restore")

	flag.StringVar(&settings.Key, "k", settings.Key, "hash Key")

	flag.StringVar(&settings.DBConnString, "d", settings.DBConnString, "db connection string")

	flag.StringVar(&settings.CryptoKey, "crypto-key", settings.CryptoKey, "crypto-key")

	flag.Parse()

}

func (settings *SettingsList) setConfigEnv() {
	var cfg struct {
		Addr          string `env:"ADDRESS"`
		StoreInterval string `env:"STORE_INTERVAL"`
		StoreFile     string `env:"STORE_FILE"`
		Restore       string `env:"RESTORE"`
		Key           string `env:"KEY"`
		DBConnString  string `env:"DATABASE_DSN"`
		CryptoKey     string `env:"CRYPTO_KEY"`
	}

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if len(strings.TrimSpace(cfg.Addr)) != 0 {
		settings.Addr = cfg.Addr
	}

	if len(strings.TrimSpace(cfg.StoreInterval)) != 0 {
		settings.StoreInterval, err = time.ParseDuration(cfg.StoreInterval)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.Restore)) != 0 {
		settings.Restore, err = strconv.ParseBool(cfg.Restore)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.Key)) != 0 {
		settings.Key = cfg.Key
	}

	if len(strings.TrimSpace(cfg.DBConnString)) != 0 {
		settings.DBConnString = cfg.DBConnString
	}

	if len(strings.TrimSpace(cfg.CryptoKey)) != 0 {
		settings.CryptoKey = cfg.CryptoKey
	}

}
