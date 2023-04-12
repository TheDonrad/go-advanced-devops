package main

import (
	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/agent/config"
	"goAdvancedTpl/internal/agent/sender"
	"log"
	"runtime"
	"time"
)

func main() {

	settings := config.SetConfig()

	metrics := collector.NewMetrics()
	var memStats runtime.MemStats
	startTime := time.Now()
	for {
		metrics.SetMetrics(memStats)

		if time.Since(startTime) >= settings.ReportInterval {
			metrics.CalculateMetrics()
			err := sender.SendMetrics(settings.Addr, metrics, settings.Key)
			if err != nil {
				log.Println(err.Error())
			}
			metrics.SetMetricsToZero()
			startTime = time.Now()

		}
		<-time.After(settings.PollInterval)
	}

}
