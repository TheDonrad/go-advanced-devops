package main

import (
	"fmt"
	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/agent/sender"
	"runtime"
	"time"
)

func main() {

	settings := setConfig()

	metrics := collector.NewMetrics()
	var memStats runtime.MemStats
	startTime := time.Now()
	for {
		metrics.SetMetrics(memStats)

		if time.Since(startTime) >= settings.reportInterval {
			metrics.CalculateMetrics()
			err := sender.SendMetrics(settings.addr, metrics, settings.key)
			if err != nil {
				fmt.Println(err.Error())
			}
			metrics.SetMetricsToZero()
			startTime = time.Now()

		}
		<-time.After(settings.pollInterval)
	}

}
