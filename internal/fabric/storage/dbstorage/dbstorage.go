package dbstorage

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"

	"goAdvancedTpl/internal/fabric/metrics"

	"github.com/jackc/pgx/v5"
)

type DBStorage struct {
	Metrics  metrics.MetricStorage
	Settings struct {
		DBConnString string
	}
	Mutex sync.RWMutex
}

// NewDBStorage создаёт объект для хранения метрик в СУБД
func NewDBStorage(connString string, restore bool) *DBStorage {
	s := &DBStorage{
		Metrics:  metrics.NewMetricStorage(),
		Settings: struct{ DBConnString string }{DBConnString: connString},
	}
	if restore {
		s.Restore()
	}
	return s
}

func (m *DBStorage) AddValue(metricType string, metricName string, f float64, i int64) {
	m.Mutex.Lock()
	switch metricType {
	case "gauge":
		m.Metrics.Gauge[metricName] = f
	default:
		m.Metrics.Counter[metricName] += i
	}
	m.Mutex.Unlock()
	if err := m.Save(); err != nil {
		return
	}
}

func (m *DBStorage) Ping() (err error) {

	if len(m.Settings.DBConnString) == 0 {
		return nil
	}

	conn, err := pgx.Connect(context.Background(), m.Settings.DBConnString)
	if err != nil {
		return err
	}

	defer func() {
		if err = conn.Close(context.Background()); err != nil {
			log.Print(err)
		}
	}()

	if err = conn.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}

func (m *DBStorage) Render(w http.ResponseWriter) (err error) {

	return m.Metrics.Render(w)
}

func (m *DBStorage) GetIntValue(metricName string) (value int64, err error) {
	m.Mutex.RLock()
	value, ok := m.Metrics.Counter[metricName]
	m.Mutex.RUnlock()
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *DBStorage) GetFloatValue(metricName string) (value float64, err error) {
	m.Mutex.RLock()
	value, ok := m.Metrics.Gauge[metricName]
	m.Mutex.RUnlock()
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *DBStorage) GetValue(metricType string, metricName string) (str string, err error) {

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
