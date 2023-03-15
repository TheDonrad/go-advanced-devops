package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type (
	Gauge   float64
	Counter int64
)

type MemStorage struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	PollCount     Counter
	RandomValue   Gauge
}

var Metrics MemStorage

func (m *MemStorage) getStringValue(metricType string, metricName string) (string, error) {

	var err error
	var value string
	switch metricType {
	case "Gauge":
		switch metricName {
		case "Alloc":
			value = gaugeToStr(Metrics.Alloc)
		case "BuckHashSys":
			value = gaugeToStr(Metrics.BuckHashSys)
		case "Frees":
			value = gaugeToStr(Metrics.Frees)
		case "GCCPUFraction":
			value = gaugeToStr(Metrics.GCCPUFraction)
		case "GCSys":
			value = gaugeToStr(Metrics.GCSys)
		case "HeapAlloc":
			value = gaugeToStr(Metrics.HeapAlloc)
		case "HeapIdle":
			value = gaugeToStr(Metrics.HeapIdle)
		case "HeapInuse":
			value = gaugeToStr(Metrics.HeapInuse)
		case "HeapObjects":
			value = gaugeToStr(Metrics.HeapObjects)
		case "HeapReleased":
			value = gaugeToStr(Metrics.HeapReleased)
		case "HeapSys":
			value = gaugeToStr(Metrics.HeapSys)
		case "LastGC":
			value = gaugeToStr(Metrics.LastGC)
		case "Lookups":
			value = gaugeToStr(Metrics.Lookups)
		case "MCacheInuse":
			value = gaugeToStr(Metrics.MCacheInuse)
		case "MCacheSys":
			value = gaugeToStr(Metrics.MCacheSys)
		case "MSpanInuse":
			value = gaugeToStr(Metrics.MSpanInuse)
		case "MSpanSys":
			value = gaugeToStr(Metrics.MSpanSys)
		case "Mallocs":
			value = gaugeToStr(Metrics.Mallocs)
		case "NextGC":
			value = gaugeToStr(Metrics.NextGC)
		case "NumForcedGC":
			value = gaugeToStr(Metrics.NumForcedGC)
		case "NumGC":
			value = gaugeToStr(Metrics.NumGC)
		case "OtherSys":
			value = gaugeToStr(Metrics.OtherSys)
		case "PauseTotalNs":
			value = gaugeToStr(Metrics.PauseTotalNs)
		case "StackInuse":
			value = gaugeToStr(Metrics.StackInuse)
		case "StackSys":
			value = gaugeToStr(Metrics.StackSys)
		case "Sys":
			value = gaugeToStr(Metrics.Sys)
		case "TotalAlloc":
			value = gaugeToStr(Metrics.TotalAlloc)
		case "RandomValue":
			value = gaugeToStr(Metrics.RandomValue)
		case "PollCount":
			value = counterToStr(Metrics.PollCount)
		default:
			err = errors.New("not found")
		}
	case "Counter":
		switch metricName {
		case "PollCount":
			value = counterToStr(Metrics.PollCount)
		default:
			err = errors.New("not found")
		}
	}
	return value, err
}

func gaugeToStr(f Gauge) string {
	return strconv.FormatFloat(float64(f), 'g', 1, 64)
}

func counterToStr(f Counter) string {
	return strconv.FormatInt(int64(f), 2)
}

func Main() {

	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", writeMetric)
	r.Get("/value/{metricType}/{metricName}", GetMetric)
	r.Get("/", AllMetrics)
	er := http.ListenAndServe(":8080", r)
	if er != nil {
		fmt.Println(er.Error())
	}
}
