package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"html/template"
	"net/http"
	"strconv"
)

type MetricStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMetricStorage() *MetricStorage {
	return &MetricStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (m *MetricStorage) AddGauge(metricName string, value float64, dbConnString string) {
	m.Gauge[metricName] = value
	if len(dbConnString) > 0 {
		if err := m.Save(dbConnString, ""); err != nil {
			fmt.Println(err.Error())
		}
	}

}

func (m *MetricStorage) AddCounter(metricName string, value int64, dbConnString string) {
	m.Counter[metricName] += value
	if len(dbConnString) > 0 {
		if err := m.Save(dbConnString, ""); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (m *MetricStorage) GetIntValue(metricName string) (value int64, err error) {
	value, ok := m.Counter[metricName]
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *MetricStorage) GetFloatValue(metricName string) (value float64, err error) {
	value, ok := m.Gauge[metricName]
	if !ok {
		err = errors.New("no such metric")
	}
	return
}

func (m *MetricStorage) GetValue(metricType string, metricName string) (str string, err error) {

	switch metricType {
	case "gauge":
		value, ok := m.Gauge[metricName]
		if !ok {
			err = errors.New("no such metric")
			return
		}
		str = strconv.FormatFloat(value, 'f', -1, 64)
	case "counter":
		value, ok := m.Counter[metricName]
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

func (m *MetricStorage) Render(w http.ResponseWriter) error {
	content := pageTemplate()

	tmpl, err := template.New("metrics_page").Parse(content)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, m)
	return err
}

func pageTemplate() string {
	content := `
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<title>Metrics</title>
		</head>
		<body>
		<table>
			<thead>
				<tr>
					<th>metric</th>
					<th>value</th>
				</tr>
			</thead>
			<tbody>
				{{range $metric, $value := .Gauge }}
				<tr>
					<td>{{ $metric }}</td>
					<td>{{ $value }}</td>
				</tr>
				{{ end }}
				{{range $metric, $value := .Counter }}
				<tr>
					<td>{{ $metric }}</td>
					<td>{{ $value }}</td>
				</tr>
				{{ end }}
			</tbody>
			</tbody>
		</table>			
		</body>
		</html>`
	return content
}

func (m *MetricStorage) Ping(dbConnString string) (err error) {

	if len(dbConnString) == 0 {
		return nil
	}

	conn, err := pgx.Connect(context.Background(), dbConnString)
	if err != nil {
		return err
	}

	defer func() {
		if err = conn.Close(context.Background()); err != nil {
			fmt.Println(err)
		}
	}()

	if err = conn.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}
