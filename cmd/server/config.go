package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"strconv"
	"strings"
	"time"
)

type serverConfig struct {
	addr          string
	storeInterval time.Duration
	storeFile     string
	restore       bool
	key           string
	dbConnString  string
}

func srvConfig() serverConfig {
	srvConfig := serverConfig{
		addr:          "127.0.0.1:8080",
		storeInterval: 300 * time.Second,
		storeFile:     "/tmp/devops-metrics-db.json",
		restore:       true,
		dbConnString:  "",
	}
	srvConfig.setConfigFlags()
	srvConfig.setConfigEnv()

	return srvConfig
}

func (srvConfig *serverConfig) setConfigFlags() {

	flag.StringVar(&srvConfig.addr, "a", srvConfig.addr, "host to listen on")

	flag.StringVar(&srvConfig.storeFile, "f", srvConfig.storeFile, "file to store metrics")

	flag.Func("i", "store interval", func(flagValue string) error {
		srvConfig.storeInterval, _ = time.ParseDuration(flagValue)
		return nil
	})
	flag.BoolVar(&srvConfig.restore, "r", srvConfig.restore, "restore")

	flag.StringVar(&srvConfig.key, "k", srvConfig.key, "hash key")

	flag.StringVar(&srvConfig.dbConnString, "d", srvConfig.dbConnString, "db connection string")

	flag.Parse()

}

func (srvConfig *serverConfig) setConfigEnv() {
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
		fmt.Println(err.Error())
		return
	}

	if len(strings.TrimSpace(cfg.Addr)) != 0 {
		srvConfig.addr = cfg.Addr
	}

	if len(strings.TrimSpace(cfg.StoreInterval)) != 0 {
		srvConfig.storeInterval, err = time.ParseDuration(cfg.StoreInterval)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.Restore)) != 0 {
		srvConfig.restore, err = strconv.ParseBool(cfg.Restore)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.Key)) != 0 {
		srvConfig.key = cfg.Key
	}

	if len(strings.TrimSpace(cfg.DBConnString)) != 0 {
		srvConfig.dbConnString = cfg.DBConnString
	}
}
