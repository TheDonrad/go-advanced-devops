package main

import (
	"encoding/json"
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
	GetIntValue(metricType string, metricName string) (value int64, err error)
	GetFloatValue(metricType string, metricName string) (value float64, err error)
}

type Metric struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
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

func (h *APIHandler) WriteWholeMetric(w http.ResponseWriter, r *http.Request) {

	met := Metric{}
	if err := json.NewDecoder(r.Body).Decode(&met); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch met.MType {
	case "gauge":
		h.metrics.AddGauge(met.ID, met.Value)
	case "counter":
		h.metrics.AddCounter(met.ID, met.Delta)
	default:
		http.Error(w, "Invalid metric type", http.StatusNotImplemented)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(met)
	_, err := w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func (h *APIHandler) GetWholeMetric(w http.ResponseWriter, r *http.Request) {

	var err error
	met := Metric{}
	if err = json.NewDecoder(r.Body).Decode(&met); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch met.MType {
	case "gauge":
		met.Value, err = h.metrics.GetFloatValue(met.MType, met.ID)
	case "counter":
		met.Delta, err = h.metrics.GetIntValue(met.MType, met.ID)
	default:
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(met)
	_, err = w.Write(b)
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
