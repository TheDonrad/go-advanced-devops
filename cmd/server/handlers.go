package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
)

func writeMetric(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	switch metricType {
	case "gauge":
		valueFloat, _ := strconv.ParseFloat(metricValue, 64)
		value := Gauge(valueFloat)
		switch metricName {
		case "Alloc":
			Metrics.Alloc = value
		case "BuckHashSys":
			Metrics.BuckHashSys = value
		case "Frees":
			Metrics.Frees = value
		case "GCCPUFraction":
			Metrics.GCCPUFraction = value
		case "GCSys":
			Metrics.GCSys = value
		case "HeapAlloc":
			Metrics.HeapAlloc = value
		case "HeapIdle":
			Metrics.HeapIdle = value
		case "HeapInuse":
			Metrics.HeapInuse = value
		case "HeapObjects":
			Metrics.HeapObjects = value
		case "HeapReleased":
			Metrics.HeapReleased = value
		case "HeapSys":
			Metrics.HeapSys = value
		case "LastGC":
			Metrics.LastGC = value
		case "Lookups":
			Metrics.Lookups = value
		case "MCacheInuse":
			Metrics.MCacheInuse = value
		case "MCacheSys":
			Metrics.MCacheSys = value
		case "MSpanInuse":
			Metrics.MSpanInuse = value
		case "MSpanSys":
			Metrics.MSpanSys = value
		case "Mallocs":
			Metrics.Mallocs = value
		case "NextGC":
			Metrics.NextGC = value
		case "NumForcedGC":
			Metrics.NumForcedGC = value
		case "NumGC":
			Metrics.NumGC = value
		case "OtherSys":
			Metrics.OtherSys = value
		case "PauseTotalNs":
			Metrics.PauseTotalNs = value
		case "StackInuse":
			Metrics.StackInuse = value
		case "StackSys":
			Metrics.StackSys = value
		case "Sys":
			Metrics.Sys = value
		case "TotalAlloc":
			Metrics.TotalAlloc = value
		case "RandomValue":
			Metrics.RandomValue = value
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	//setValue(metricName, value)
	case "counter":
		valueInt, _ := strconv.ParseInt(metricValue, 0, 8)
		value := Counter(valueInt)
		switch metricName {
		case "PollCount":
			Metrics.PollCount = value
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	//setValue(metricName, value)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue, err := Metrics.getStringValue(metricType, metricName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = w.Write([]byte(metricValue))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func AllMetrics(w http.ResponseWriter, _ *http.Request) {
	content := pageTemplate()

	tmpl, err := template.New("metrics_page").Parse(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "Metrics", Metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func pageTemplate() string {
	content := `
		{{define "Metrics"}}
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<title>Metrics</title>
		</head>
		<body>
			<table class="table">
			<th>metric</th>
			<th>value</th>
			<tr><td>Alloc</td><td>{{.Alloc}}</td></tr>
			<tr><td>BuckHashSys</td><td>{{.BuckHashSys}}</td></tr>
			<tr><td>Frees</td><td>{{.Frees}}</td></tr>
			<tr><td>GCCPUFraction</td><td>{{.GCCPUFraction}}</td></tr>
			<tr><td>GCSys</td><td>{{.GCSys}}</td></tr>
			<tr><td>HeapAlloc</td><td>{{.HeapAlloc}}</td></tr>
			<tr><td>HeapIdle</td><td>{{.HeapIdle}}</td></tr>
			<tr><td>HeapInuse</td><td>{{.HeapInuse}}</td></tr>
			<tr><td>HeapObjects</td><td>{{.HeapObjects}}</td></tr>
			<tr><td>HeapReleased</td><td>{{.HeapReleased}}</td></tr>
			<tr><td>HeapSys</td><td>{{.HeapSys}}</td></tr>
			<tr><td>LastGC</td><td>{{.LastGC}}</td></tr>
			<tr><td>Lookups</td><td>{{.Lookups}}</td></tr>
			<tr><td>MCacheInuse</td><td>{{.MCacheInuse}}</td></tr>
			<tr><td>MCacheSys</td><td>{{.MCacheSys}}</td></tr>
			<tr><td>MSpanInuse</td><td>{{.MSpanInuse}}</td></tr>
			<tr><td>MSpanSys</td><td>{{.MSpanSys}}</td></tr>
			<tr><td>Mallocs</td><td>{{.Mallocs}}</td></tr>
			<tr><td>NextGC</td><td>{{.NextGC}}</td></tr>
			<tr><td>NumForcedGC</td><td>{{.NumForcedGC}}</td></tr>
			<tr><td>NumGC</td><td>{{.NumGC}}</td></tr>
			<tr><td>OtherSys</td><td>{{.OtherSys}}</td></tr>
			<tr><td>PauseTotalNs</td><td>{{.PauseTotalNs}}</td></tr>
			<tr><td>StackInuse</td><td>{{.StackInuse}}</td></tr>
			<tr><td>StackSys</td><td>{{.StackSys}}</td></tr>
			<tr><td>Sys</td><td>{{.Sys}}</td></tr>
			<tr><td>TotalAlloc</td><td>{{.TotalAlloc}}</td></tr>
			<tr><td>RandomValue</td><td>{{.RandomValue}}</td></tr>
			{{end}}
		</table>
		</body>
		</html>`
	return content
}

//func setValue(metricName string, value Gauge) {
//	// pointer to struct - addressable
//	ps := reflect.ValueOf(&Metrics)
//	// struct
//	s := ps.Elem()
//	if s.Kind() == reflect.Struct {
//		// exported field
//		f := s.FieldByName(metricName)
//		if f.IsValid() {
//			// A Value can be changed only if it is
//			// addressable and was not obtained by
//			// the use of unexported struct fields.
//			if f.CanSet() {
//				//f.SetInt(value)
//			}
//		}
//	}
//}
