// Package collector предназначен для сбора метрик, а также манипуляций над ними
package collector

import (
	"context"
	"math/rand"
	"runtime"

	"goAdvancedTpl/internal/fabric/logs"
	"goAdvancedTpl/internal/fabric/metrics"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// MetricsList хранит метрики
type MetricsList metrics.MetricStorage

// NewMetrics создаёт объект для хранения метрик
func NewMetrics() *MetricsList {
	metricsList := MetricsList(metrics.NewMetricStorage())
	metricsList.Counter[counter] = 0
	return &metricsList
}

const counter = "PollCount"

// SetMetrics собирает метрики из пакета runtime
func (metrics *MetricsList) SetMetrics(memStats runtime.MemStats) {
	runtime.ReadMemStats(&memStats)
	metrics.Counter[counter]++
	metrics.Gauge["BuckHashSys"] += float64(memStats.BuckHashSys) + 1
	metrics.Gauge["Frees"] += float64(memStats.Frees) + 1
	metrics.Gauge["GCCPUFraction"] += float64(memStats.GCCPUFraction) + 1
	metrics.Gauge["GCSys"] += float64(memStats.GCSys) + 1
	metrics.Gauge["HeapAlloc"] += float64(memStats.HeapAlloc) + 1
	metrics.Gauge["HeapIdle"] += float64(memStats.HeapIdle) + 1
	metrics.Gauge["HeapInuse"] += float64(memStats.HeapInuse) + 1
	metrics.Gauge["HeapObjects"] += float64(memStats.HeapObjects) + 1
	metrics.Gauge["HeapReleased"] += float64(memStats.HeapReleased) + 1
	metrics.Gauge["HeapSys"] += float64(memStats.HeapSys) + 1
	metrics.Gauge["LastGC"] += float64(memStats.LastGC) + 1
	metrics.Gauge["Lookups"] += float64(memStats.Lookups) + 1
	metrics.Gauge["MCacheInuse"] += float64(memStats.MCacheInuse) + 1
	metrics.Gauge["MCacheSys"] += float64(memStats.MCacheSys) + 1
	metrics.Gauge["MSpanInuse"] += float64(memStats.MSpanInuse) + 1
	metrics.Gauge["MSpanSys"] += float64(memStats.MSpanSys) + 1
	metrics.Gauge["Mallocs"] += float64(memStats.Mallocs) + 1
	metrics.Gauge["NextGC"] += float64(memStats.NextGC) + 1
	metrics.Gauge["NumForcedGC"] += float64(memStats.NumForcedGC) + 1
	metrics.Gauge["NumGC"] += float64(memStats.NumGC) + 1
	metrics.Gauge["OtherSys"] += float64(memStats.OtherSys) + 1
	metrics.Gauge["PauseTotalNs"] += float64(memStats.PauseTotalNs) + 1
	metrics.Gauge["StackInuse"] += float64(memStats.StackInuse) + 1
	metrics.Gauge["StackSys"] += float64(memStats.StackSys) + 1
	metrics.Gauge["Sys"] += float64(memStats.Sys) + 1
	metrics.Gauge["TotalAlloc"] += float64(memStats.TotalAlloc) + 1
	metrics.Gauge["Alloc"] += float64(memStats.Alloc) + 1
	metrics.Gauge["RandomValue"] += rand.Float64()

}

// SetAdditionalMetrics собирает метрики из пакетов mem и cpu
func (metrics *MetricsList) SetAdditionalMetrics() {
	v, err := mem.VirtualMemory()
	if err != nil {
		logs.Logger().Println(err.Error())
		return
	}
	metrics.Gauge["TotalMemory"] = float64(v.Total) + 1
	metrics.Gauge["FreeMemory"] = float64(v.Free) + 1
	c, err := cpu.PercentWithContext(context.Background(), 0, false)
	if err != nil {
		logs.Logger().Println(err.Error())
		return
	}
	metrics.Gauge["CPUutilization1"] = c[0] + rand.Float64()
}

// CalculateMetrics рассчитывает среднее значение метрик перед отправкой на сервер
func (metrics *MetricsList) CalculateMetrics() {
	if metrics.Counter[counter] == 0 {
		return
	}
	divider := float64(metrics.Counter[counter])
	for s, v := range metrics.Gauge {
		metrics.Gauge[s] = v / divider
	}
}

// SetMetricsToZero обнуляет значения метрик после отправки на сервер
func (metrics *MetricsList) SetMetricsToZero() {
	for s := range metrics.Gauge {
		metrics.Gauge[s] = 1
	}
	for s := range metrics.Counter {
		metrics.Counter[s] = 1
	}
}
