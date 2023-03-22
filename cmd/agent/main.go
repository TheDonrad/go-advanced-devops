package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/agent/sender"
	"runtime"
	"strings"
	"time"
)

type config struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval string `env:"REPORT_INTERVAL"`
	PollInterval   string `env:"POLL_INTERVAL"`
}

func main() {

	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	addr := "127.0.0.1:8080"
	if len(strings.TrimSpace(cfg.Addr)) != 0 {
		addr = cfg.Addr
	}

	pollInterval := 2 * time.Second
	if len(strings.TrimSpace(cfg.PollInterval)) != 0 {
		pollInterval, err = time.ParseDuration(cfg.PollInterval)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	reportInterval := 5 * time.Second
	if len(strings.TrimSpace(cfg.ReportInterval)) != 0 {
		reportInterval, err = time.ParseDuration(cfg.ReportInterval)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	metrics := collector.NewMetrics()
	var memStats runtime.MemStats
	startTime := time.Now()
	for {
		metrics.SetMetrics(memStats)

		if time.Since(startTime) >= reportInterval {
			metrics.CalculateMetrics()
			err := sender.SendMetrics(addr, metrics)
			if err != nil {
				fmt.Println(err.Error())
			}
			metrics.SetMetricsToZero()
			startTime = time.Now()

		}
		<-time.After(pollInterval)
	}

}
