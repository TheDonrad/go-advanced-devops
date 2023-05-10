// сервис для сбора метрик ОС
package main

import (

	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/agent/config"
	"goAdvancedTpl/internal/agent/sender"
	"goAdvancedTpl/internal/fabric/logs"
	"goAdvancedTpl/internal/fabric/onstart"
)

var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
)

func main() {

	onstart.WriteMessage(BuildVersion, BuildDate, BuildCommit)

	settings := config.Config(true)

	metrics := collector.NewMetrics()
	var memStats runtime.MemStats
	wg := &sync.WaitGroup{}
	var sendingInProgress int32
	wg.Add(1)
	go func() {
		for {
			if atomic.LoadInt32(&sendingInProgress) == 1 {
				time.Sleep(time.Second)
				continue
			}
			metrics.SetMetrics(memStats)
			time.Sleep(settings.ReportInterval)

		}
	}()

	wg.Add(1)
	go func() {
		for {
			if atomic.LoadInt32(&sendingInProgress) == 1 {
				time.Sleep(time.Second)
				continue
			}

			metrics.SetAdditionalMetrics()
			time.Sleep(settings.ReportInterval)
		}
	}()

	wg.Add(1)
	go func() {

		for {
			time.Sleep(settings.PollInterval)
			atomic.StoreInt32(&sendingInProgress, 1)
			metrics.CalculateMetrics()

			err := sender.SendMetrics(settings.Addr, metrics, settings.Key, settings.RateLimit)
			if err != nil {
				logs.New().Println(err.Error())
			}

			metrics.SetMetricsToZero()

			if err != nil {
				logs.New().Println(err.Error())
			}
			atomic.StoreInt32(&sendingInProgress, 0)
		}
	}()

	wg.Wait()
}
