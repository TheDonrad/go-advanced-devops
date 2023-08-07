package grpcserver

import (
	"context"

	"goAdvancedTpl/internal/fabric/calchash"
	"goAdvancedTpl/internal/fabric/proto"
)

type IStorage interface {
	AddValue(metricType string, metricName string, f float64, i int64)
}

type MetricServer struct {
	proto.UnimplementedMetricsServer
	metrics IStorage
	key     string
}

func NewServer(met IStorage, key string) *MetricServer {
	return &MetricServer{
		metrics: met,
		key:     key,
	}
}

func (s *MetricServer) AddMetrics(_ context.Context, in *proto.AddMetricsRequest) (*proto.AddMetricsResponse, error) {
	var response proto.AddMetricsResponse

	for _, m := range in.Metric {
		var hash string
		if m.MType == "gauge" {
			hash = calchash.Calculate[float64](s.key, "gauge", m.ID, m.Value)
		} else {
			hash = calchash.Calculate[int64](s.key, "counter", m.ID, m.Delta)
		}
		if hash != m.Hash {
			response.Error = "hash is incorrect"
			return &response, nil
		}
		s.metrics.AddValue(m.MType, m.ID, m.Value, m.Delta)
	}
	return &response, nil
}
