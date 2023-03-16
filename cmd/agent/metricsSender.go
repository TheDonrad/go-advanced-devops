package main

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func SendMetrics(metrics *metricsList) error {

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
		switch typeOf.Elem().Field(i).Type.String() {
		case "main.gauge":
			t = "gauge"
		default:
			t = "counter"
		}
		endpoint := fmt.Sprintf("http://%s/update/%s/%s/%s",
			"127.0.0.1:8080", t, typeOf.Elem().Field(i).Name, v)
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
	}
	return nil
}
