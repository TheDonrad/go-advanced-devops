package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"goAdvancedTpl/internal/fabric/calchash"

	"github.com/go-chi/chi/v5"
)

type IStorage interface {
	AddValue(metricType string, metricName string, f float64, i int64)
	Render(w http.ResponseWriter) error
	GetValue(metricType string, metricName string) (string, error)
	GetIntValue(metricName string) (value int64, err error)
	GetFloatValue(metricName string) (value float64, err error)
	Ping() (err error)
	Save() (err error)
}

type Metric struct {
	ID    string `json:"id"`              // имя метрики
	MType string `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64  `json:"delta,omitempty"` // значение метрики в случае передачи counter
	// TODO: проверить указатель и без omitempty
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string  `json:"hash,omitempty"`  // значение хеш-функции
}

type APIHandler struct {
	metrics IStorage
	key     string
}

func NewAPIHandler(metrics IStorage, key string) (h *APIHandler) {
	h = &APIHandler{
		metrics: metrics,
		key:     key,
	}
	return
}

// WriteMetric записывает метрику переданную в адресе HTTP-запроса
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
		h.metrics.AddValue(metricType, metricName, valueFloat, 0)
	case "counter":
		valueInt, err := strconv.ParseInt(metricValue, 0, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.metrics.AddValue(metricType, metricName, 0, valueInt)
	default:
		http.Error(w, "Invalid metric type", http.StatusNotImplemented)
	}
	w.WriteHeader(http.StatusOK)
}

// WriteWholeMetric записывает метрику Metric, переданную в формате JSON
func (h *APIHandler) WriteWholeMetric(w http.ResponseWriter, r *http.Request) {

	met := Metric{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&met); err != nil && !errors.Is(err, io.EOF) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch met.MType {
	case "gauge":
		h.metrics.AddValue(met.MType, met.ID, met.Value, 0)
		hash := met.Hash
		met.Hash = calchash.Calculate[float64](h.key, met.MType, met.ID, met.Value)
		if hash != met.Hash {
			http.Error(w, "Invalid hash", http.StatusBadRequest)
		}
	case "counter":
		h.metrics.AddValue(met.MType, met.ID, 0, met.Delta)
		hash := met.Hash
		met.Hash = calchash.Calculate[int64](h.key, met.MType, met.ID, met.Delta)
		if hash != met.Hash {
			http.Error(w, "Invalid hash", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Invalid metric type", http.StatusNotImplemented)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(met)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetMetric возвращает значение метрики строкой
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

// GetWholeMetric возвращает метрику Metric в формате JSON
func (h *APIHandler) GetWholeMetric(w http.ResponseWriter, r *http.Request) {

	var err error
	met := Metric{}
	if err = json.NewDecoder(r.Body).Decode(&met); err != nil && !errors.Is(err, io.EOF) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: вынести в отдельный слой
	switch met.MType {
	case "gauge":
		met.Value, err = h.metrics.GetFloatValue(met.ID)
		met.Hash = calchash.Calculate[float64](h.key, met.MType, met.ID, met.Value)
	case "counter":
		met.Delta, err = h.metrics.GetIntValue(met.ID)
		met.Hash = calchash.Calculate[int64](h.key, met.MType, met.ID, met.Delta)
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

// AllMetrics возвращает все метрики в формате HTML
func (h *APIHandler) AllMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	err := h.metrics.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Ping проверяет доступность хранилища
func (h *APIHandler) Ping(w http.ResponseWriter, _ *http.Request) {

	err := h.metrics.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(""))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// WriteAllMetrics записывает массив метрик Metric, переданный в формате JSON
func (h *APIHandler) WriteAllMetrics(w http.ResponseWriter, r *http.Request) {

	var met []Metric
	var err error
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &met); err != nil && !errors.Is(err, io.EOF) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, met := range met {
		switch met.MType {
		case "gauge":
			h.metrics.AddValue(met.MType, met.ID, met.Value, 0)
			hash := met.Hash
			met.Hash = calchash.Calculate[float64](h.key, met.MType, met.ID, met.Value)
			if hash != met.Hash {
				http.Error(w, "Invalid hash", http.StatusBadRequest)
			}
		case "counter":
			h.metrics.AddValue(met.MType, met.ID, 0, met.Delta)
			hash := met.Hash
			met.Hash = calchash.Calculate[int64](h.key, met.MType, met.ID, met.Delta)
			if hash != met.Hash {
				http.Error(w, "Invalid hash", http.StatusBadRequest)
			}

		default:
			http.Error(w, "Invalid metric type", http.StatusNotImplemented)
		}
	}

	b, _ := json.Marshal(met) // Для обхода ошибки автотестов
	if err = h.metrics.Save(); err != nil {
		log.Println(err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
