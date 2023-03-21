package main

import (
	"goAdvancedTpl/internal/agent/collector"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateMetrics(t *testing.T) {
	tests := []struct {
		name    string
		metrics collector.MetricsList
		want    collector.Gauge
	}{
		{
			name:    "calc",
			want:    collector.Gauge(5),
			metrics: collector.MetricsList{Alloc: collector.Gauge(10), PollCount: collector.Counter(2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metrics.CalculateMetrics()
			assert.Equal(t, tt.want, tt.metrics.Alloc)
		})
	}
}
