package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

type Metric struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string  `json:"hash,omitempty"`  // значение хеш-функции
}

func SendMetrics(addr string, metrics interface{}, key string) (err error) {

	client := &http.Client{}
	values := reflect.ValueOf(metrics)
	typeOf := reflect.TypeOf(metrics)

	for i := 0; i < typeOf.Elem().NumField(); i++ {
		var v string
		value := values.Elem().Field(i)
		switch typeOf.Elem().Field(i).Type.Kind() {
		case reflect.Float64:
			v = strconv.FormatFloat(value.Float(), 'g', 4, 64)
		case reflect.Int64:
			v = strconv.FormatInt(value.Int(), 10)
		}
		var t string
		var met Metric
		switch typeOf.Elem().Field(i).Type.String() {
		case "collector.Gauge":
			t = "gauge"
			met = Metric{
				ID:    typeOf.Elem().Field(i).Name,
				MType: t,
				Value: value.Float(),
			}
			calcGaugeHash(&met, key)
		default:
			t = "counter"
			met = Metric{
				ID:    typeOf.Elem().Field(i).Name,
				MType: t,
				Delta: value.Int(),
			}
			calcCounterHash(&met, key)
		}

		endpoint := fmt.Sprintf("http://%s/update/%s/%s/%s",
			addr, t, typeOf.Elem().Field(i).Name, v)
		request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(""))
		if err != nil {
			return err
		}
		request.Header.Add("Content-Type", "text/plain")
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		err = response.Body.Close()
		if err != nil {
			return err
		}
		endpoint = fmt.Sprintf("http://%s/update/", addr)
		b, _ := json.Marshal(met)
		request, err = http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(string(b)))
		if err != nil {
			return err
		}
		request.Header.Add("Content-Type", "application/json")
		response1, err := client.Do(request)
		if err != nil {
			return err
		}
		err = response1.Body.Close()
		if err != nil {
			return err
		}
	}

	return err
}
