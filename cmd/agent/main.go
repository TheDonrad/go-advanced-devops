package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/agent/sender"
	"runtime"
	"strconv"
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

	i := 2
	if len(strings.TrimSpace(cfg.PollInterval)) != 0 {
		i, err = strconv.Atoi(cfg.PollInterval)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	pollInterval := time.Duration(i) * time.Second

	j := 5
	if len(strings.TrimSpace(cfg.ReportInterval)) != 0 {
		j, err = strconv.Atoi(cfg.ReportInterval)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	reportInterval := time.Duration(j) * time.Second

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
