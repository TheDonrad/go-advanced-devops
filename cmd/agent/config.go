package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"strings"
	"time"
)

type settingsList struct {
	addr           string
	reportInterval time.Duration
	pollInterval   time.Duration
	key            string
}

func setConfig() settingsList {

	settings := settingsList{
		addr:           "127.0.0.1:8080",
		pollInterval:   2 * time.Second,
		reportInterval: 5 * time.Second,
		key:            "",
	}
	settings.setConfigFlags()
	settings.setConfigEnv()
	return settings
}

func (settings *settingsList) setConfigFlags() {

	flag.StringVar(&settings.addr, "a", settings.addr, "host to send")
	flag.Func("p", "poll interval", func(flagValue string) error {
		settings.pollInterval, _ = time.ParseDuration(flagValue)
		return nil
	})
	flag.Func("r", "report interval", func(flagValue string) error {
		settings.reportInterval, _ = time.ParseDuration(flagValue)
		return nil
	})
	flag.StringVar(&settings.key, "k", settings.key, "hash key")
	flag.Parse()
}

func (settings *settingsList) setConfigEnv() {

	var cfg struct {
		Addr           string `env:"ADDRESS"`
		ReportInterval string `env:"REPORT_INTERVAL"`
		PollInterval   string `env:"POLL_INTERVAL"`
		Key            string `env:"KEY"`
	}

	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(strings.TrimSpace(cfg.Addr)) != 0 {
		settings.addr = cfg.Addr
	}

	if len(strings.TrimSpace(cfg.PollInterval)) != 0 {
		settings.pollInterval, err = time.ParseDuration(cfg.PollInterval)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.ReportInterval)) != 0 {
		settings.reportInterval, err = time.ParseDuration(cfg.ReportInterval)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.Key)) != 0 {
		settings.key = cfg.Key
	}

}
