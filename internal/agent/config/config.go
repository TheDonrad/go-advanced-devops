// Package config служит для определения настроек агента по сбору метрик
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"goAdvancedTpl/internal/fabric/logs"

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
	configFile     string        // Файл с настройками
}

// Config возвращает настройки агента из переменных окружения или флагов запуска.
// У переменных окружения приоритет перед флагами
func Config(parseFlags bool) *SettingsList {

	settings := SettingsList{}

	if parseFlags {
		settings.setConfigFlags()
	}

	if settings.configFile != "" {
		settings.setConfigFile()
	}

	settings.setConfigEnv()

	settings.setUnspecified()

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

	flag.StringVar(&settings.configFile, "c", settings.configFile, "config")
	flag.StringVar(&settings.configFile, "config", settings.configFile, "config")

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
		Config         string `env:"CONFIG"`
	}

	err := env.Parse(&cfg)
	if err != nil {
		logs.Logger().Println(err.Error())
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

	if len(strings.TrimSpace(cfg.Config)) != 0 {
		settings.configFile = cfg.Config
	}

}

func (settings *SettingsList) setUnspecified() {

	if settings.Addr == "" {
		settings.Addr = "127.0.0.1:8080"
	}

	if settings.PollInterval == 0 {
		settings.PollInterval = 2 * time.Second
	}

	if settings.ReportInterval == 0 {
		settings.ReportInterval = 5 * time.Second
	}

	if settings.RateLimit == 0 {
		settings.RateLimit = 5
	}

}

func (settings *SettingsList) setConfigFile() {

	file, err := os.ReadFile(settings.configFile)
	if err != nil {
		logs.Logger().Println(err.Error())
		return
	}

	var cfg struct {
		Addr           string `json:"address"`
		ReportInterval string `json:"report_interval"`
		PollInterval   string `json:"poll_interval"`
		CryptoKey      string `json:"crypto_key"`
	}

	if err = json.Unmarshal(file, &cfg); err != nil {
		logs.Logger().Println(err.Error())
		return
	}

	if settings.Addr == "" {
		settings.Addr = cfg.Addr
	}

	if settings.PollInterval == 0 {
		settings.PollInterval, err = time.ParseDuration(cfg.PollInterval)
		if err != nil {
			logs.Logger().Println(err.Error())
		}
	}

	if settings.ReportInterval == 0 {
		settings.ReportInterval, err = time.ParseDuration(cfg.ReportInterval)
		if err != nil {
			logs.Logger().Println(err.Error())
		}
	}

	if settings.CryptoKey == "" {
		settings.CryptoKey = cfg.CryptoKey
	}

}
