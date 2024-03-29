// Package config служит для определения настроек серевера сбора метрик
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"goAdvancedTpl/internal/fabric/logs"

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
	configFile    string        // Файл с настройками
}

// Config возвращает настройки агента из переменных окружения или флагов запуска.
// У переменных окружения приоритет перед флагами
func Config() *SettingsList {

	settings := SettingsList{}

	settings.setConfigFlags()

	if settings.configFile != "" {
		settings.setConfigFile()
	}

	settings.setConfigEnv()

	settings.setUnspecified()

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

	flag.StringVar(&settings.configFile, "c", settings.configFile, "config")
	flag.StringVar(&settings.configFile, "config", settings.configFile, "config")

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
		Config        string `env:"CONFIG"`
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

	if len(strings.TrimSpace(cfg.Config)) != 0 {
		settings.configFile = cfg.Config
	}

}

func (settings *SettingsList) setUnspecified() {

	if settings.Addr == "" {
		settings.Addr = "127.0.0.1:8080"
	}

	if settings.StoreInterval == 0 {
		settings.StoreInterval = 120 * time.Second
	}

	if settings.StoreFile == "" {
		settings.StoreFile = "/tmp/devops-metrics-db.json"
	}

}

func (settings *SettingsList) setConfigFile() {

	file, err := os.ReadFile(settings.configFile)
	if err != nil {
		logs.Logger().Println(err.Error())
		return
	}

	var cfg struct {
		Addr          string `json:"address"`
		Restore       bool   `json:"restore"`
		StoreInterval string `json:"store_interval"`
		StoreFile     string `json:"store_file"`
		DBConnString  string `json:"database_dsn"`
		CryptoKey     string `json:"crypto_key"`
	}

	if err = json.Unmarshal(file, &cfg); err != nil {
		logs.Logger().Println(err.Error())
		return
	}

	if settings.Addr == "" {
		settings.Addr = cfg.Addr
	}

	if !settings.Restore {
		settings.Restore = cfg.Restore
	}

	if settings.StoreInterval == 0 {
		settings.StoreInterval, err = time.ParseDuration(cfg.StoreInterval)
		if err != nil {
			logs.Logger().Println(err.Error())
		}
	}

	if settings.StoreFile == "" {
		settings.StoreFile = cfg.StoreFile
	}

	if settings.DBConnString == "" {
		settings.DBConnString = cfg.DBConnString
	}

	if settings.CryptoKey == "" {
		settings.CryptoKey = cfg.CryptoKey
	}

}
