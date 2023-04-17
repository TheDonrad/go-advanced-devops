package filestorage

import (
	"errors"
	"goAdvancedTpl/internal/fabric/metricsstorage"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FileStorage struct {
	Metrics  metricsstorage.MetricStorage
	Settings struct {
		StoreFile     string
		StoreInterval time.Duration
	}
}

func NewFileStorage(storeInterval time.Duration, storeFile string) *FileStorage {
	fs := &FileStorage{
		Metrics: metricsstorage.NewMetricStorage(),
		Settings: struct {
			StoreFile     string
			StoreInterval time.Duration
		}{StoreFile: storeFile, StoreInterval: storeInterval},
	}
	go func() {
		for {
			<-time.After(storeInterval)
			if err := fs.Save(); err != nil {
				log.Println(err.Error())
			}
		}
	}()
	return fs
}

func (m *FileStorage) AddValue(metricType string, metricName string, f float64, i int64) {
	switch metricType {
	case "gauge":
		m.Metrics.Gauge[metricName] = f
	default:
		m.Metrics.Counter[metricName] += i
	}
}

func (m *FileStorage) Ping() (err error) {
	return nil
}

func (m *FileStorage) Render(w http.ResponseWriter) (err error) {

	return m.Metrics.Render(w)
}

func (m *FileStorage) GetIntValue(metricName string) (value int64, err error) {
	value, ok := m.Metrics.Counter[metricName]
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *FileStorage) GetFloatValue(metricName string) (value float64, err error) {
	value, ok := m.Metrics.Gauge[metricName]
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *FileStorage) GetValue(metricType string, metricName string) (str string, err error) {

	switch metricType {
	case "gauge":
		value, ok := m.Metrics.Gauge[metricName]
		if !ok {
			err = errors.New("no such metric")
			return
		}
		str = strconv.FormatFloat(value, 'f', -1, 64)
	case "counter":
		value, ok := m.Metrics.Counter[metricName]
		if !ok {
			err = errors.New("no such metric")
			return
		}
		str = strconv.FormatInt(value, 10)
	default:
		err = errors.New("no such metric")
	}

	return
}
