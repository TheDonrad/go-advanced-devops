package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"strings"
	"time"
)

type SettingsList struct {
	Addr           string
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string
}

func SetConfig() *SettingsList {

	settings := SettingsList{
		Addr:           "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 5 * time.Second,
		Key:            "",
	}
	settings.setConfigFlags()
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
	flag.Parse()
}

func (settings *SettingsList) setConfigEnv() {

	var cfg struct {
		Addr           string `env:"ADDRESS"`
		ReportInterval string `env:"REPORT_INTERVAL"`
		PollInterval   string `env:"POLL_INTERVAL"`
		Key            string `env:"KEY"`
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

}