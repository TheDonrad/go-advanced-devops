package main

import (
	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/agent/config"
	"goAdvancedTpl/internal/agent/sender"
	"goAdvancedTpl/internal/fabric/logs"
	"runtime"
	"sync"
	"time"
)

func main() {

	settings := config.SetConfig()

	metrics := collector.NewMetrics()
	var memStats runtime.MemStats
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {

		for {
			metrics.SetMetrics(memStats)
			time.Sleep(settings.ReportInterval)
		}
	}()
	go func() {
		wg.Add(1)
		for {
			metrics.SetAdditionalMetrics()
			time.Sleep(settings.ReportInterval)
		}
	}()
	wg.Add(1)
	go func() {

		for {

			time.Sleep(settings.PollInterval)
			metrics.CalculateMetrics()
			err := sender.SendMetrics(settings.Addr, metrics, settings.Key)
			if err != nil {
				logs.New().Println(err.Error())
			}
			metrics.SetMetricsToZero()

		}
	}()

	wg.Wait()
}
