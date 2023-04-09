package metricsstorage

type MetricStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMetricStorage() MetricStorage {
	return MetricStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}
