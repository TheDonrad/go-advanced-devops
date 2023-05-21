// сервис для сбора метрик ОС
package main

import (
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

	idleConnClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-sigint
		close(idleConnClosed)
	}()

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

			<-idleConnClosed
			wg.Done()
			break
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

			<-idleConnClosed
			wg.Done()
			break
		}
	}()

	wg.Add(1)
	go func() {
	Loop:
		for {
			select {
			case <-idleConnClosed:
				wg.Done()
				break Loop
			default:
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
