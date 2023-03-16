package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type metricsHandlers interface {
	AddGauge(metricName string, value float64)
	AddCounter(metricName string, value int64)
	GetValue(metricType string, metricName string) (string, error)
	Render(w http.ResponseWriter) error
}

type APIHandler struct {
	metrics metricsHandlers
}

func NewAPIHandler(metric metricsHandlers) (h *APIHandler) {
	h = &APIHandler{metrics: metric}
	return
}
func (h *APIHandler) WriteMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricType = strings.ToLower(metricType)
	metricName := chi.URLParam(r, "metricName")
	metricType = strings.ToLower(metricType)
	metricValue := chi.URLParam(r, "metricValue")
	switch metricType {
	case "gauge":
		valueFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.metrics.AddGauge(metricName, valueFloat)
	case "counter":
		valueInt, err := strconv.ParseInt(metricValue, 0, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.metrics.AddCounter(metricName, valueInt)
	default:
		http.Error(w, "Invalid metric type", http.StatusNotImplemented)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) GetMetric(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricType = strings.ToLower(metricType)
	metricName := chi.URLParam(r, "metricName")
	metricType = strings.ToLower(metricType)

	metricValue, err := h.metrics.GetValue(metricType, metricName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(metricValue))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *APIHandler) AllMetrics(w http.ResponseWriter, _ *http.Request) {
	err := h.metrics.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
