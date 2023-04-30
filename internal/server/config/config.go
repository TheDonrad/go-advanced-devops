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

// ServerConfig хранит настройки сервера
// Если заполнен параметр DBConnString, считаем, что в качестве хранилища следует использовать БД
// Иначе - файл
type ServerConfig struct {
	Addr          string        // Адрес для получения метрик
	StoreInterval time.Duration // Период сохранения настроек
	StoreFile     string        // Путь к файлу для хранения метрик
	Restore       bool          // Восстанавливать метрики из хранилища при запуске
	Key           string        // Ключ для отправки шифрованного хеша метрики по алгоритму sha256
	DBConnString  string        // Строка соединения с БД
}

// SrvConfig возвращает настройки агента из переменных окружения или флагов запуска.
// У переменных окружения приоритет перед флагами
func SrvConfig() *ServerConfig {
	srvConfig := ServerConfig{
		Addr:          "127.0.0.1:8080",
		StoreInterval: 120 * time.Second,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       true,
		DBConnString:  "",
	}
	srvConfig.setConfigFlags()
	srvConfig.setConfigEnv()

	return &srvConfig
}

func (srvConfig *ServerConfig) setConfigFlags() {

	flag.StringVar(&srvConfig.Addr, "a", srvConfig.Addr, "host to listen on")

	flag.StringVar(&srvConfig.StoreFile, "f", srvConfig.StoreFile, "file to store metrics")

	flag.Func("i", "store interval", func(flagValue string) error {
		srvConfig.StoreInterval, _ = time.ParseDuration(flagValue)
		return nil
	})
	flag.BoolVar(&srvConfig.Restore, "r", srvConfig.Restore, "restore")

	flag.StringVar(&srvConfig.Key, "k", srvConfig.Key, "hash Key")

	flag.StringVar(&srvConfig.DBConnString, "d", srvConfig.DBConnString, "db connection string")

	flag.Parse()

}

func (srvConfig *ServerConfig) setConfigEnv() {
	var cfg struct {
		Addr          string `env:"ADDRESS"`
		StoreInterval string `env:"STORE_INTERVAL"`
		StoreFile     string `env:"STORE_FILE"`
		Restore       string `env:"RESTORE"`
		Key           string `env:"KEY"`
		DBConnString  string `env:"DATABASE_DSN"`
	}

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if len(strings.TrimSpace(cfg.Addr)) != 0 {
		srvConfig.Addr = cfg.Addr
	}

	if len(strings.TrimSpace(cfg.StoreInterval)) != 0 {
		srvConfig.StoreInterval, err = time.ParseDuration(cfg.StoreInterval)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.Restore)) != 0 {
		srvConfig.Restore, err = strconv.ParseBool(cfg.Restore)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.Key)) != 0 {
		srvConfig.Key = cfg.Key
	}

	if len(strings.TrimSpace(cfg.DBConnString)) != 0 {
		srvConfig.DBConnString = cfg.DBConnString
	}
}
