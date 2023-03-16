package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {

	pollInterval := 2 * time.Second
	reportInterval := 5 * time.Second

	metrics := NewMetrics()
	var memStats runtime.MemStats
	startTime := time.Now()
	for {
		metrics.SetMetrics(memStats)

		if time.Since(startTime) >= reportInterval {
			metrics.CalculateMetrics()
			err := SendMetrics(metrics)
			if err != nil {
				fmt.Println(err.Error())
			}
			metrics.SetMetricsToZero()
			startTime = time.Now()

		}
		<-time.After(pollInterval)
	}

}
