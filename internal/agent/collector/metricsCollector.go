package collector

import (
	"math/rand"
	"runtime"
)

type (
	Gauge   float64
	Counter int64
)

type MetricsList struct {
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

func NewMetrics() *MetricsList {
	return &MetricsList{}
}

func (metrics *MetricsList) SetMetrics(memStats runtime.MemStats) {
	runtime.ReadMemStats(&memStats)
	metrics.BuckHashSys += Gauge(memStats.BuckHashSys) + 1
	metrics.Frees += Gauge(memStats.Frees) + 1
	metrics.GCCPUFraction += Gauge(memStats.GCCPUFraction) + 1
	metrics.GCSys += Gauge(memStats.GCSys) + 1
	metrics.HeapAlloc += Gauge(memStats.HeapAlloc) + 1
	metrics.HeapIdle += Gauge(memStats.HeapIdle) + 1
	metrics.HeapInuse += Gauge(memStats.HeapInuse) + 1
	metrics.HeapObjects += Gauge(memStats.HeapObjects) + 1
	metrics.HeapReleased += Gauge(memStats.HeapReleased) + 1
	metrics.HeapSys += Gauge(memStats.HeapSys) + 1
	metrics.LastGC += Gauge(memStats.LastGC) + 1
	metrics.Lookups += Gauge(memStats.Lookups) + 1
	metrics.MCacheInuse += Gauge(memStats.MCacheInuse) + 1
	metrics.MCacheSys += Gauge(memStats.MCacheSys) + 1
	metrics.MSpanInuse += Gauge(memStats.MSpanInuse) + 1
	metrics.MSpanSys += Gauge(memStats.MSpanSys) + 1
	metrics.Mallocs += Gauge(memStats.Mallocs) + 1
	metrics.NextGC += Gauge(memStats.NextGC) + 1
	metrics.NumForcedGC += Gauge(memStats.NumForcedGC) + 1
	metrics.NumGC += Gauge(memStats.NumGC) + 1
	metrics.OtherSys += Gauge(memStats.OtherSys) + 1
	metrics.PauseTotalNs += Gauge(memStats.PauseTotalNs) + 1
	metrics.StackInuse += Gauge(memStats.StackInuse) + 1
	metrics.StackSys += Gauge(memStats.StackSys) + 1
	metrics.Sys += Gauge(memStats.Sys) + 1
	metrics.TotalAlloc += Gauge(memStats.TotalAlloc) + 1
	metrics.RandomValue += Gauge(rand.Float64())
	metrics.Alloc += Gauge(memStats.Alloc)
	metrics.PollCount++
}

func (metrics *MetricsList) CalculateMetrics() {
	metrics.Alloc = metrics.Alloc / Gauge(metrics.PollCount)
	metrics.BuckHashSys = metrics.BuckHashSys / Gauge(metrics.PollCount)
	metrics.Frees = metrics.Frees / Gauge(metrics.PollCount)
	metrics.GCCPUFraction = metrics.GCCPUFraction / Gauge(metrics.PollCount)
	metrics.GCSys = metrics.GCSys / Gauge(metrics.PollCount)
	metrics.HeapAlloc = metrics.HeapAlloc / Gauge(metrics.PollCount)
	metrics.HeapIdle = metrics.HeapIdle / Gauge(metrics.PollCount)
	metrics.HeapInuse = metrics.HeapInuse / Gauge(metrics.PollCount)
	metrics.HeapObjects = metrics.HeapObjects / Gauge(metrics.PollCount)
	metrics.HeapReleased = metrics.HeapReleased / Gauge(metrics.PollCount)
	metrics.HeapSys = metrics.HeapSys / Gauge(metrics.PollCount)
	metrics.LastGC = metrics.LastGC / Gauge(metrics.PollCount)
	metrics.Lookups = metrics.Lookups / Gauge(metrics.PollCount)
	metrics.MCacheInuse = metrics.MCacheInuse / Gauge(metrics.PollCount)
	metrics.MCacheSys = metrics.MCacheSys / Gauge(metrics.PollCount)
	metrics.MSpanInuse = metrics.MSpanInuse / Gauge(metrics.PollCount)
	metrics.MSpanSys = metrics.MSpanSys / Gauge(metrics.PollCount)
	metrics.Mallocs = metrics.Mallocs / Gauge(metrics.PollCount)
	metrics.NextGC = metrics.NextGC / Gauge(metrics.PollCount)
	metrics.NumForcedGC = metrics.NumForcedGC / Gauge(metrics.PollCount)
	metrics.NumGC = metrics.NumGC / Gauge(metrics.PollCount)
	metrics.OtherSys = metrics.OtherSys / Gauge(metrics.PollCount)
	metrics.PauseTotalNs = metrics.PauseTotalNs / Gauge(metrics.PollCount)
	metrics.StackInuse = metrics.StackInuse / Gauge(metrics.PollCount)
	metrics.StackSys = metrics.StackSys / Gauge(metrics.PollCount)
	metrics.Sys = metrics.Sys / Gauge(metrics.PollCount)
	metrics.TotalAlloc = metrics.TotalAlloc / Gauge(metrics.PollCount)
	metrics.RandomValue = metrics.RandomValue / Gauge(metrics.PollCount)
}

func (metrics *MetricsList) SetMetricsToZero() {
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
