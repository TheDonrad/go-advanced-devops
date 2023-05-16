package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/fabric/calchash"
	"goAdvancedTpl/internal/fabric/encryption"

	"golang.org/x/sync/errgroup"
)

// Metric служит для сериализации значения метрики в json
type Metric struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Hash  string  `json:"hash,omitempty"`  // значение хеш-функции
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type metricsMap []Metric

// SendMetrics отправляет метрики на сервер
func SendMetrics(addr string, metrics *collector.MetricsList, key string, limit int, cryptoKey string) (err error) {

	client := &http.Client{}
	length := len(metrics.Gauge) + len(metrics.Counter)
	metricsToSend := make(metricsMap, length)
	i := 0
	for name, value := range metrics.Gauge {

		met := Metric{
			ID:    name,
			MType: "gauge",
			Value: value,
		}
		met.Hash = calchash.Calculate[float64](key, met.MType, met.ID, met.Value)
		metricsToSend[i] = met
		i++

	}

	for name, value := range metrics.Counter {
		met := Metric{
			ID:    name,
			MType: "counter",
			Delta: value,
		}
		met.Hash = calchash.Calculate[int64](key, met.MType, met.ID, met.Delta)
		metricsToSend[i] = met
		i++
	}

	if limit <= 1 {

		err = sendCollection(addr, metricsToSend, client)
		if err != nil {
			return err
		}

		if err = sendAllMetrics(addr, metricsToSend, client, cryptoKey); err != nil {
			return err
		}
	} else {
		g := &errgroup.Group{}
		for j := 0; j < len(metricsToSend)-1; j++ {
			g.Go(func() error {
				lErr := sendCollection(addr, metricsToSend[j:j+1], client)
				if lErr != nil {
					return lErr
				}
				return nil
			})
			if j > 0 && j%limit == 0 {
				time.Sleep(time.Second)
			}
		}

		g.Go(func() error {
			lErr := sendAllMetrics(addr, metricsToSend, client, cryptoKey)
			if lErr != nil {
				return lErr
			}
			return nil
		})
		if err = g.Wait(); err != nil {
			return err
		}
	}
	return err
}

func sendCollection(addr string, metricsToSend metricsMap, client *http.Client) (err error) {
	for _, m := range metricsToSend {
		if err = sendOneString(addr, m, client); err != nil {
			return err
		}
		if err = sendJSON(addr, m, client); err != nil {
			return err
		}
	}
	return nil
}

func sendOneString(addr string, met Metric, client *http.Client) error {
	var endpoint string
	if met.MType == "gauge" {
		endpoint = "http://" + addr + "/update/" + met.MType + "/" + met.ID + "/" +
			strconv.FormatFloat(met.Value, 'E', -1, 64)
	} else {
		endpoint = "http://" + addr + "/update/" + met.MType + "/" + met.ID + "/" +
			strconv.FormatInt(met.Delta, 10)
	}
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if err = response.Body.Close(); err != nil {
		return err
	}
	return nil
}

func sendJSON(addr string, met Metric, client *http.Client) error {
	endpoint := fmt.Sprintf("http://%s/update/", addr)
	b, _ := json.Marshal(met)
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(string(b)))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if err = response.Body.Close(); err != nil {
		return err
	}
	return nil
}

func sendAllMetrics(addr string, met metricsMap, client *http.Client, cryptoKey string) error {
	endpoint := fmt.Sprintf("http://%s/updates/", addr)
	b, _ := json.Marshal(met)

	if cryptoKey != "" {
		res, err := encryption.Encrypt(cryptoKey, b)
		if err == nil {
			b = res
		}
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(string(b)))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if err = response.Body.Close(); err != nil {
		return err
	}
	return nil
}
