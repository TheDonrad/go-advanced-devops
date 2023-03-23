package main

import (
	"fmt"
	"net/http"
)

type (
	gauge   float64
	counter int64
)

type MemStorage struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
}

var metrics MemStorage

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/update/", writeMetric)
	// запуск сервера с адресом localhost, порт 8080
	er := http.ListenAndServe(":8080", nil)
	if er != nil {
		fmt.Println(er.Error())
	}
}
