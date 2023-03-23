package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
			calculateMetrics(&tt.metrics)
			assert.Equal(t, tt.want, tt.metrics.Alloc)
		})
	}
}
