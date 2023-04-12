package main

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/servermiddleware"
	"goAdvancedTpl/internal/server/storage"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {

	srvConfig := srvConfig()

	metStorage := storage.NewMetricStorage()
	savingSettings := storage.NewSavingSettings(srvConfig.storeInterval, srvConfig.storeFile)
	metStorage.Restore(srvConfig.restore, savingSettings, metStorage)
	go func() {
		for {
			<-time.After(srvConfig.storeInterval)
			metStorage.Save(savingSettings)
		}
	}()
	r := routers(metStorage, srvConfig.key)
	er := http.ListenAndServe(srvConfig.addr, r)
	if er != nil {
		fmt.Println(er.Error())
	}

}

func routers(metStorage *storage.MetricStorage, key string) *chi.Mux {

	h := handlers.NewAPIHandler(metStorage, key)
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(servermiddleware.GzipHandle)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.WriteWholeMetric)
		r.Post("/{metricType}/{metricName}/{metricValue}", h.WriteMetric)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetWholeMetric)
		r.Get("/{metricType}/{metricName}", h.GetMetric)
	})
	r.Get("/", h.AllMetrics)
	return r
}
