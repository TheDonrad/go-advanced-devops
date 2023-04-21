package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goAdvancedTpl/internal/agent/collector"
	"goAdvancedTpl/internal/fabric/calchash"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type Metric struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string  `json:"hash,omitempty"`  // значение хеш-функции
}

type metricsMap []Metric

func SendMetrics(addr string, metrics *collector.MetricsList, key string, limit int) (err error) {

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
		met.Hash = calchash.Calculate(key, met.MType, met.ID, met.Value)
		metricsToSend[i] = met
		i++

	}

	for name, value := range metrics.Counter {
		met := Metric{
			ID:    name,
			MType: "counter",
			Delta: value,
		}
		met.Hash = calchash.Calculate(key, met.MType, met.ID, met.Delta)
		metricsToSend[i] = met
		i++
	}

	if limit <= 1 {

		err = sendCollection(addr, metricsToSend, client)
		if err != nil {
			return err
		}

		if err = sendAllMetrics(addr, metricsToSend, client); err != nil {
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
			lErr := sendAllMetrics(addr, metricsToSend, client)
			if lErr != nil {
				return lErr
			}
			return nil
		})
		if err := g.Wait(); err != nil {
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
		endpoint = fmt.Sprintf("http://%s/update/%s/%s/%f",
			addr, met.MType, met.ID, met.Value)
	} else {
		endpoint = fmt.Sprintf("http://%s/update/%s/%s/%d",
			addr, met.MType, met.ID, met.Delta)
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

func sendAllMetrics(addr string, met metricsMap, client *http.Client) error {
	endpoint := fmt.Sprintf("http://%s/updates/", addr)
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
