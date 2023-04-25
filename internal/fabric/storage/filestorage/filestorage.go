package filestorage

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"goAdvancedTpl/internal/fabric/metrics"
)

type FileStorage struct {
	Metrics  metrics.MetricStorage
	Settings struct {
		StoreFile     string
		StoreInterval time.Duration
	}
	Mutex sync.RWMutex
}

func NewFileStorage(storeInterval time.Duration, storeFile string) *FileStorage {
	fs := &FileStorage{
		Metrics: metrics.NewMetricStorage(),
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
	m.Mutex.Lock()
	switch metricType {
	case "gauge":
		m.Metrics.Gauge[metricName] = f
	default:
		m.Metrics.Counter[metricName] += i
	}
	m.Mutex.Unlock()
}

func (m *FileStorage) Ping() (err error) {
	return nil
}

func (m *FileStorage) Render(w http.ResponseWriter) (err error) {

	return m.Metrics.Render(w)
}

func (m *FileStorage) GetIntValue(metricName string) (value int64, err error) {
	m.Mutex.RLock()
	value, ok := m.Metrics.Counter[metricName]
	m.Mutex.RUnlock()
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *FileStorage) GetFloatValue(metricName string) (value float64, err error) {
	m.Mutex.RLock()
	value, ok := m.Metrics.Gauge[metricName]
	m.Mutex.RUnlock()
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *FileStorage) GetValue(metricType string, metricName string) (str string, err error) {

	switch metricType {
	case "gauge":
		m.Mutex.RLock()
		value, ok := m.Metrics.Gauge[metricName]
		m.Mutex.RUnlock()
		if !ok {
			err = errors.New("no such metric")
			return
		}
		str = strconv.FormatFloat(value, 'f', -1, 64)
	case "counter":
		m.Mutex.RLock()
		value, ok := m.Metrics.Counter[metricName]
		m.Mutex.RUnlock()
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
