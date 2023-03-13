package main

import (
	"net/http"
	"strconv"
	"strings"
)

func writeMetric(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.RequestURI, "/")
	if len(path) != 5 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metricType := path[2]
	metricName := path[3]
	metricValue := path[4]
	switch metricType {
	case "gauge":
		valueFloat, _ := strconv.ParseFloat(metricValue, 64)
		value := gauge(valueFloat)
		switch metricName {
		case "Alloc":
			metrics.Alloc = value
		case "BuckHashSys":
			metrics.BuckHashSys = value
		case "Frees":
			metrics.Frees = value
		case "GCCPUFraction":
			metrics.GCCPUFraction = value
		case "GCSys":
			metrics.GCSys = value
		case "HeapAlloc":
			metrics.HeapAlloc = value
		case "HeapIdle":
			metrics.HeapIdle = value
		case "HeapInuse":
			metrics.HeapInuse = value
		case "HeapObjects":
			metrics.HeapObjects = value
		case "HeapReleased":
			metrics.HeapReleased = value
		case "HeapSys":
			metrics.HeapSys = value
		case "LastGC":
			metrics.LastGC = value
		case "Lookups":
			metrics.Lookups = value
		case "MCacheInuse":
			metrics.MCacheInuse = value
		case "MCacheSys":
			metrics.MCacheSys = value
		case "MSpanInuse":
			metrics.MSpanInuse = value
		case "MSpanSys":
			metrics.MSpanSys = value
		case "Mallocs":
			metrics.Mallocs = value
		case "NextGC":
			metrics.NextGC = value
		case "NumForcedGC":
			metrics.NumForcedGC = value
		case "NumGC":
			metrics.NumGC = value
		case "OtherSys":
			metrics.OtherSys = value
		case "PauseTotalNs":
			metrics.PauseTotalNs = value
		case "StackInuse":
			metrics.StackInuse = value
		case "StackSys":
			metrics.StackSys = value
		case "Sys":
			metrics.Sys = value
		case "TotalAlloc":
			metrics.TotalAlloc = value
		case "RandomValue":
			metrics.RandomValue = value
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
		//setValue(metricName, value)
	case "counter":
		valueInt, _ := strconv.ParseInt(metricValue, 0, 8)
		value := counter(valueInt)
		switch metricName {
		case "PollCount":
			metrics.PollCount = value

		default:
			w.WriteHeader(http.StatusBadRequest)
		}
		//setValue(metricName, value)
	}

	w.WriteHeader(http.StatusOK)
}

//func setValue(metricName string, value gauge) {
//	// pointer to struct - addressable
//	ps := reflect.ValueOf(&metrics)
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
