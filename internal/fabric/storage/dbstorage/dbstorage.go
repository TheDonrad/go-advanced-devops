package dbstorage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"goAdvancedTpl/internal/fabric/metricsstorage"
	"log"
	"net/http"
	"strconv"
)

type DBStorage struct {
	Metrics  *metricsstorage.MetricStorage
	Settings struct {
		DBConnString string
	}
}

func NewDBStorage(connString string) *DBStorage {
	return &DBStorage{
		Metrics:  metricsstorage.NewMetricStorageLink(),
		Settings: struct{ DBConnString string }{DBConnString: connString},
	}
}

func (m *DBStorage) AddValue(metricType string, metricName string, f float64, i int64) {
	switch metricType {
	case "gauge":
		m.Metrics.Gauge[metricName] = f
	default:
		m.Metrics.Counter[metricName] += i
	}

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
	value, ok := m.Metrics.Counter[metricName]
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *DBStorage) GetFloatValue(metricName string) (value float64, err error) {
	value, ok := m.Metrics.Gauge[metricName]
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *DBStorage) GetValue(metricType string, metricName string) (str string, err error) {

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
