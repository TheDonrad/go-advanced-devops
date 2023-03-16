package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	metStorage := NewMetricStorage()
	h := NewAPIHandler(metStorage)
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.WriteMetric)
	r.Get("/value/{metricType}/{metricName}", h.GetMetric)
	r.Get("/", h.AllMetrics)
	er := http.ListenAndServe(":8080", r)
	if er != nil {
		fmt.Println(er.Error())
	}
}
