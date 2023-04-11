package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"goAdvancedTpl/internal/fabric/calchash"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type storage interface {
	AddGauge(metricName string, value float64, dbConnString string)
	AddCounter(metricName string, value int64, dbConnString string)
	GetValue(metricType string, metricName string) (string, error)
	Render(w http.ResponseWriter) error
	GetIntValue(metricName string) (value int64, err error)
	GetFloatValue(metricName string) (value float64, err error)
	Ping(dbConnString string) (err error)
	Save(database string, file string) (err error)
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
	metrics      storage
	key          string
	dbConnString string
}

func NewAPIHandler(metrics storage, key string, dbConnString string) (h *APIHandler) {
	h = &APIHandler{
		metrics:      metrics,
		key:          key,
		dbConnString: dbConnString,
	}
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
		h.metrics.AddGauge(metricName, valueFloat, h.dbConnString)
	case "counter":
		valueInt, err := strconv.ParseInt(metricValue, 0, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.metrics.AddCounter(metricName, valueInt, h.dbConnString)
	default:
		http.Error(w, "Invalid metric type", http.StatusNotImplemented)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) WriteWholeMetric(w http.ResponseWriter, r *http.Request) {

	met := Metric{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&met); err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch met.MType {
	case "gauge":
		h.metrics.AddGauge(met.ID, met.Value, h.dbConnString)
		hash := met.Hash
		met.Hash = calchash.Calculate(h.key, met.MType, met.ID, met.Value)
		if hash != met.Hash {
			http.Error(w, "Invalid hash", http.StatusBadRequest)
		}
	case "counter":
		h.metrics.AddCounter(met.ID, met.Delta, h.dbConnString)
		hash := met.Hash
		met.Hash = calchash.Calculate(h.key, met.MType, met.ID, met.Delta)
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
	if err = json.NewDecoder(r.Body).Decode(&met); err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: вынести в отдельный слой
	switch met.MType {
	case "gauge":
		met.Value, err = h.metrics.GetFloatValue(met.ID)
		met.Hash = calchash.Calculate(h.key, met.MType, met.ID, met.Value)
	case "counter":
		met.Delta, err = h.metrics.GetIntValue(met.ID)
		met.Hash = calchash.Calculate(h.key, met.MType, met.ID, met.Delta)
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
	w.Header().Add("Content-Type", "text/html")
	err := h.metrics.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *APIHandler) Ping(w http.ResponseWriter, _ *http.Request) {

	err := h.metrics.Ping(h.dbConnString)
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

func (h *APIHandler) WriteAllMetrics(w http.ResponseWriter, r *http.Request) {

	var met []Metric
	var err error
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &met); err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	type received struct {
		gauge map[string]Metric
		count map[string]Metric
	}

	addedValues := received{
		make(map[string]Metric),
		make(map[string]Metric),
	}
	for _, met := range met {
		switch met.MType {
		case "gauge":
			h.metrics.AddGauge(met.ID, met.Value, "")
			hash := met.Hash
			met.Hash = calchash.Calculate(h.key, met.MType, met.ID, met.Value)
			if hash != met.Hash {
				err = errors.New("invalid hash")
				http.Error(w, "Invalid hash", http.StatusBadRequest)
			}
			addedValues.gauge[met.ID] = met
		case "counter":
			h.metrics.AddCounter(met.ID, met.Delta, "")
			hash := met.Hash
			met.Hash = calchash.Calculate(h.key, met.MType, met.ID, met.Delta)
			if hash != met.Hash {
				err = errors.New("invalid hash")
				http.Error(w, "Invalid hash", http.StatusBadRequest)
			}
			addedValues.count[met.ID] = met
		default:
			http.Error(w, "Invalid metric type", http.StatusNotImplemented)
		}
	}

	var sendMet []Metric
	for _, metric := range addedValues.gauge {
		sendMet = append(sendMet, metric)
	}
	for _, metric := range addedValues.count {
		sendMet = append(sendMet, metric)
	}
	b, _ := json.Marshal(sendMet[0])
	if err = h.metrics.Save(h.dbConnString, ""); err != nil {
		fmt.Println(err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
