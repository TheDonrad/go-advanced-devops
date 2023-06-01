// сервис для сбора метрик ОС
package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-sigint
		cancel()
	}()

	metrics := collector.NewMetrics()
	var memStats runtime.MemStats
	wg := &sync.WaitGroup{}
	var sendingInProgress int32
	wg.Add(1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				if atomic.LoadInt32(&sendingInProgress) == 1 {
					time.Sleep(time.Second)
					continue
				}
				metrics.SetMetrics(memStats)
				time.Sleep(settings.ReportInterval)

			}
		}
	}()

	wg.Add(1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				if atomic.LoadInt32(&sendingInProgress) == 1 {
					time.Sleep(time.Second)
					continue
				}

				metrics.SetAdditionalMetrics()
				time.Sleep(settings.ReportInterval)

			}
		}
	}()

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(settings.PollInterval)
	Loop:
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				wg.Done()
				break Loop
			case <-ticker.C:
				time.Sleep(settings.PollInterval)
				atomic.StoreInt32(&sendingInProgress, 1)
				metrics.CalculateMetrics()

				err := sender.SendMetrics(settings.Addr, metrics, settings.Key, settings.RateLimit, settings.CryptoKey)
				if err != nil {
					logs.Logger().Println(err.Error())
				}

				metrics.SetMetricsToZero()

				if err != nil {
					logs.Logger().Println(err.Error())
				}
				atomic.StoreInt32(&sendingInProgress, 0)
			}
		}
	}()

	wg.Wait()
}
