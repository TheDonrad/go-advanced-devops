package main

import (
	"fmt"
	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/agent/sender"
	"runtime"
	"time"
)

func main() {

	pollInterval := 2 * time.Second
	reportInterval := 5 * time.Second

	metrics := collector.NewMetrics()
	var memStats runtime.MemStats
	startTime := time.Now()
	for {
		metrics.SetMetrics(memStats)

		if time.Since(startTime) >= reportInterval {
			metrics.CalculateMetrics()
			err := sender.SendMetrics(metrics)
			if err != nil {
				fmt.Println(err.Error())
			}
			metrics.SetMetricsToZero()
			startTime = time.Now()

		}
		<-time.After(pollInterval)
	}

}
