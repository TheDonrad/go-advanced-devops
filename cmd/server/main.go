package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/storage"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type config struct {
	Addr string `env:"ADDRESS"`
}

func main() {

	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println(err.Error())
	}
	addr := "127.0.0.1:8080"
	if len(strings.TrimSpace(cfg.Addr)) != 0 {
		addr = cfg.Addr
	}

	metStorage := storage.NewMetricStorage()
	h := handlers.NewAPIHandler(metStorage)
	r := chi.NewRouter()

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.WriteWholeMetric)
		r.Post("/{metricType}/{metricName}/{metricValue}", h.WriteMetric)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetWholeMetric)
		r.Get("/{metricType}/{metricName}", h.GetMetric)
	})
	r.Get("/", h.AllMetrics)
	er := http.ListenAndServe(addr, r)
	if er != nil {
		fmt.Println(er.Error())
	}
}
