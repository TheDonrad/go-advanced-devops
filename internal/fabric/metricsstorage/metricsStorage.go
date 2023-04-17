package metricsstorage

import (
	"html/template"
	"net/http"
)

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

func NewMetricStorageLink() *MetricStorage {
	return &MetricStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
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
