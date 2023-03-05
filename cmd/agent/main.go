package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

type (
	gauge   float64
	counter int64
)

type metricsList struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
}

func main() {

	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	var metrics metricsList
	var memStats runtime.MemStats
	startTime := time.Now()
	for {
		setMetrics(&metrics, memStats)
		fmt.Println("met")

		if time.Duration(time.Now().Sub(startTime)) >= reportInterval {
			fmt.Println("send")
			calculateMetrics(&metrics)
			sendMetrics(metrics)
			setMetricsToZero(&metrics)
			startTime = time.Now()

		}
		<-time.After(pollInterval)
	}

}

func setMetrics(metrics *metricsList, memStats runtime.MemStats) {
	runtime.ReadMemStats(&memStats)
	metrics.BuckHashSys += gauge(memStats.BuckHashSys)
	metrics.Frees += gauge(memStats.Frees)
	metrics.GCCPUFraction += gauge(memStats.GCCPUFraction)
	metrics.GCSys += gauge(memStats.GCSys)
	metrics.HeapAlloc += gauge(memStats.HeapAlloc)
	metrics.HeapIdle += gauge(memStats.HeapIdle)
	metrics.HeapInuse += gauge(memStats.HeapInuse)
	metrics.HeapObjects += gauge(memStats.HeapObjects)
	metrics.HeapReleased += gauge(memStats.HeapReleased)
	metrics.HeapSys += gauge(memStats.HeapSys)
	metrics.LastGC += gauge(memStats.LastGC)
	metrics.Lookups += gauge(memStats.Lookups)
	metrics.MCacheInuse += gauge(memStats.MCacheInuse)
	metrics.MCacheSys += gauge(memStats.MCacheSys)
	metrics.MSpanInuse += gauge(memStats.MSpanInuse)
	metrics.MSpanSys += gauge(memStats.MSpanSys)
	metrics.Mallocs += gauge(memStats.Mallocs)
	metrics.NextGC += gauge(memStats.NextGC)
	metrics.NumForcedGC += gauge(memStats.NumForcedGC)
	metrics.NumGC += gauge(memStats.NumGC)
	metrics.OtherSys += gauge(memStats.OtherSys)
	metrics.PauseTotalNs += gauge(memStats.PauseTotalNs)
	metrics.StackInuse += gauge(memStats.StackInuse)
	metrics.StackSys += gauge(memStats.StackSys)
	metrics.Sys += gauge(memStats.Sys)
	metrics.TotalAlloc += gauge(memStats.TotalAlloc)
	metrics.RandomValue += gauge(rand.Float64())
	metrics.PollCount++
}

func calculateMetrics(metrics *metricsList) {
	metrics.Alloc = metrics.Alloc / gauge(metrics.PollCount)
	metrics.BuckHashSys = metrics.BuckHashSys / gauge(metrics.PollCount)
	metrics.Frees = metrics.Frees / gauge(metrics.PollCount)
	metrics.GCCPUFraction = metrics.GCCPUFraction / gauge(metrics.PollCount)
	metrics.GCSys = metrics.GCSys / gauge(metrics.PollCount)
	metrics.HeapAlloc = metrics.HeapAlloc / gauge(metrics.PollCount)
	metrics.HeapIdle = metrics.HeapIdle / gauge(metrics.PollCount)
	metrics.HeapInuse = metrics.HeapInuse / gauge(metrics.PollCount)
	metrics.HeapObjects = metrics.HeapObjects / gauge(metrics.PollCount)
	metrics.HeapReleased = metrics.HeapReleased / gauge(metrics.PollCount)
	metrics.HeapSys = metrics.HeapSys / gauge(metrics.PollCount)
	metrics.LastGC = metrics.LastGC / gauge(metrics.PollCount)
	metrics.Lookups = metrics.Lookups / gauge(metrics.PollCount)
	metrics.MCacheInuse = metrics.MCacheInuse / gauge(metrics.PollCount)
	metrics.MCacheSys = metrics.MCacheSys / gauge(metrics.PollCount)
	metrics.MSpanInuse = metrics.MSpanInuse / gauge(metrics.PollCount)
	metrics.MSpanSys = metrics.MSpanSys / gauge(metrics.PollCount)
	metrics.Mallocs = metrics.Mallocs / gauge(metrics.PollCount)
	metrics.NextGC = metrics.NextGC / gauge(metrics.PollCount)
	metrics.NumForcedGC = metrics.NumForcedGC / gauge(metrics.PollCount)
	metrics.NumGC = metrics.NumGC / gauge(metrics.PollCount)
	metrics.OtherSys = metrics.OtherSys / gauge(metrics.PollCount)
	metrics.PauseTotalNs = metrics.PauseTotalNs / gauge(metrics.PollCount)
	metrics.StackInuse = metrics.StackInuse / gauge(metrics.PollCount)
	metrics.StackSys = metrics.StackSys / gauge(metrics.PollCount)
	metrics.Sys = metrics.Sys / gauge(metrics.PollCount)
	metrics.TotalAlloc = metrics.TotalAlloc / gauge(metrics.PollCount)
	metrics.RandomValue = metrics.RandomValue / gauge(metrics.PollCount)

}

func sendMetrics(metrics metricsList) {

	client := &http.Client{}
	values := reflect.ValueOf(metrics)
	types := values.Type()

	for i := 0; i < values.NumField(); i++ {
		var v string
		fmt.Println(types.Field(i).Name, values.Field(i))
		value := values.Field(i)
		switch value.Kind() {
		case reflect.Float64:
			v = strconv.FormatFloat(value.Float(), 'g', 4, 64)
		case reflect.Int64:
			v = strconv.FormatInt(value.Int(), 10)
		}
		endpoint := fmt.Sprintf("http://%s/update/%s/%s/%s",
			"127.0.0.1:8080", types.Field(i).Type, types.Field(i).Name, v)
		fmt.Println(endpoint)
		request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(""))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		request.Header.Add("Content-Type", "text/plain")
		_, err = client.Do(request)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}
}

func setMetricsToZero(metrics *metricsList) {
	metrics.BuckHashSys = 0
	metrics.Frees = 0
	metrics.GCCPUFraction = 0
	metrics.GCSys = 0
	metrics.HeapAlloc = 0
	metrics.HeapIdle = 0
	metrics.HeapInuse = 0
	metrics.HeapObjects = 0
	metrics.HeapReleased = 0
	metrics.HeapSys = 0
	metrics.LastGC = 0
	metrics.Lookups = 0
	metrics.MCacheInuse = 0
	metrics.MCacheSys = 0
	metrics.MSpanInuse = 0
	metrics.MSpanSys = 0
	metrics.Mallocs = 0
	metrics.NextGC = 0
	metrics.NumForcedGC = 0
	metrics.NumGC = 0
	metrics.OtherSys = 0
	metrics.PauseTotalNs = 0
	metrics.StackInuse = 0
	metrics.StackSys = 0
	metrics.Sys = 0
	metrics.TotalAlloc = 0
	metrics.RandomValue = 0
	metrics.PollCount = 0
}
