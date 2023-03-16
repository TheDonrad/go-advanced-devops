package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateMetrics(t *testing.T) {
	tests := []struct {
		name    string
		metrics metricsList
		want    gauge
	}{
		{
			name:    "calc",
			want:    gauge(5),
			metrics: metricsList{Alloc: gauge(10), PollCount: counter(2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metrics.CalculateMetrics()
			assert.Equal(t, tt.want, tt.metrics.Alloc)
		})
	}
}
