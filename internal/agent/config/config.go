// Package config служит для определения настроек агента по сбору метрик
package config

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

// SettingsList хранит настройки агента по сбору метрик
type SettingsList struct {
	Key            string        // Ключ для отправки шифрованного хеша метрики по алгоритму sha256
	Addr           string        // Адрес для отправки метрик
	RateLimit      int           // ограничение RPS
	ReportInterval time.Duration // Период отправки
	PollInterval   time.Duration // Период сбора
	CryptoKey      string        // Ключ шифрования
}

// Config возвращает настройки агента из переменных окружения или флагов запуска.
// У переменных окружения приоритет перед флагами
func Config(parseFlags bool) *SettingsList {

	settings := SettingsList{
		Addr:           "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 5 * time.Second,
		Key:            "",
		RateLimit:      5,
		CryptoKey:      "",
	}
	if parseFlags {
		settings.setConfigFlags()
	}

	if parseFlags {
		settings.setConfigFlags()
	}

	settings.setConfigEnv()
	return &settings
}

func (settings *SettingsList) setConfigFlags() {

	flag.StringVar(&settings.Addr, "a", settings.Addr, "host to send")
	flag.Func("p", "poll interval", func(flagValue string) error {
		settings.PollInterval, _ = time.ParseDuration(flagValue)
		return nil
	})
	flag.Func("r", "report interval", func(flagValue string) error {
		settings.ReportInterval, _ = time.ParseDuration(flagValue)
		return nil
	})
	flag.StringVar(&settings.Key, "k", settings.Key, "hash key")
	flag.IntVar(&settings.RateLimit, "l", settings.RateLimit, "rate limit")
	flag.StringVar(&settings.CryptoKey, "crypto-key", settings.CryptoKey, "crypto-key")

	flag.Parse()
}

func (settings *SettingsList) setConfigEnv() {

	var cfg struct {
		Addr           string `env:"ADDRESS"`
		ReportInterval string `env:"REPORT_INTERVAL"`
		PollInterval   string `env:"POLL_INTERVAL"`
		Key            string `env:"KEY"`
		RateLimit      int    `env:"RATE_LIMIT"`
		CryptoKey      string `env:"CRYPTO_KEY"`
	}

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if len(strings.TrimSpace(cfg.Addr)) != 0 {
		settings.Addr = cfg.Addr
	}

	if len(strings.TrimSpace(cfg.PollInterval)) != 0 {
		settings.PollInterval, err = time.ParseDuration(cfg.PollInterval)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.ReportInterval)) != 0 {
		settings.ReportInterval, err = time.ParseDuration(cfg.ReportInterval)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.Key)) != 0 {
		settings.Key = cfg.Key
	}

	if cfg.RateLimit > 0 {
		settings.RateLimit = cfg.RateLimit
	}

	if len(strings.TrimSpace(cfg.CryptoKey)) != 0 {
		settings.CryptoKey = cfg.CryptoKey
	}

}
