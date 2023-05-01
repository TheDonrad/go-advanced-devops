package collector

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsList_CalculateMetrics(t *testing.T) {
	tests := []struct {
		metrics *MetricsList
		name    string
		want    float64
	}{
		{
			name:    "calc",
			metrics: NewMetrics(),
			want:    5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metrics.Gauge["Alloc"] = 10
			tt.metrics.Counter["PollCount"] = 2
			tt.metrics.CalculateMetrics()
			assert.Equal(t, tt.want, tt.metrics.Gauge["Alloc"])
		})
	}
}

func TestMetricsList_SetAdditionalMetrics(t *testing.T) {
	tests := []struct {
		metrics *MetricsList
		name    string
	}{
		{
			name:    "SetAdditionalMetrics",
			metrics: NewMetrics(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metrics.SetAdditionalMetrics()
			assert.NotEmpty(t, tt.metrics.Gauge["TotalMemory"])
		})
	}
}

func TestMetricsList_SetMetrics(t *testing.T) {
	var memStats runtime.MemStats
	tests := []struct {
		metrics *MetricsList
		name    string
	}{
		{
			name:    "SetMetrics",
			metrics: NewMetrics(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metrics.SetMetrics(memStats)
			assert.NotEmpty(t, tt.metrics.Gauge["BuckHashSys"])
			assert.NotEmpty(t, tt.metrics.Counter["PollCount"])
		})
	}
}

func TestMetricsList_SetMetricsToZero(t *testing.T) {
	var memStats runtime.MemStats
	tests := []struct {
		metrics *MetricsList
		name    string
	}{
		{
			name:    "SetMetrics",
			metrics: NewMetrics(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metrics.SetAdditionalMetrics()
			tt.metrics.SetMetrics(memStats)
			tt.metrics.SetMetricsToZero()
			assert.Equal(t, float64(1), tt.metrics.Gauge["TotalMemory"])
			assert.Equal(t, int64(1), tt.metrics.Counter["PollCount"])
		})
	}
}

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "SetMetrics",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, NewMetrics())
		})
	}
}
