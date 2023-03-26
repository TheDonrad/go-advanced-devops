package main

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/serverMiddleware"
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
			<-time.After(5 * time.Second)
			metStorage.Save(savingSettings)
		}
	}()
	r := routers(metStorage)
	er := http.ListenAndServe(srvConfig.addr, r)
	if er != nil {
		fmt.Println(er.Error())
	}

}

func routers(metStorage *storage.MetricStorage) *chi.Mux {

	h := handlers.NewAPIHandler(metStorage)
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(serverMiddleware.GzipHandle)
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
