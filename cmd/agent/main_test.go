package main

import (
	"goAdvancedTpl/internal/agent/collector"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateMetrics(t *testing.T) {
	tests := []struct {
		name    string
		metrics *collector.MetricsList
		want    float64
	}{
		{
			name:    "calc",
			want:    5,
			metrics: collector.NewMetrics(),
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
