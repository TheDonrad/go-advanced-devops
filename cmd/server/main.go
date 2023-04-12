package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"goAdvancedTpl/internal/server/config"
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/servermiddleware"
	"goAdvancedTpl/internal/server/storage"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {

	srvConfig := config.SrvConfig()

	metStorage := storage.NewMetricStorage()
	savingSettings := storage.NewSavingSettings(srvConfig.StoreInterval, srvConfig.StoreFile, srvConfig.DBConnString)
	metStorage.Restore(srvConfig.Restore, savingSettings)
	go func() {
		for {
			<-time.After(srvConfig.StoreInterval)
			if err := metStorage.Save(savingSettings.Database, savingSettings.StoreFile); err != nil {
				log.Println(err.Error())
			}
		}
	}()
	r := routers(metStorage, srvConfig.Key, srvConfig.DBConnString)
	err := http.ListenAndServe(srvConfig.Addr, r)
	if err != nil {
		log.Println(err.Error())
	}

}

func routers(metStorage *storage.MetricStorage, key string, dbConnString string) *chi.Mux {

	h := handlers.NewAPIHandler(metStorage, key, dbConnString)
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(servermiddleware.GzipHandle)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.WriteWholeMetric)
		r.Post("/{metricType}/{metricName}/{metricValue}", h.WriteMetric)
	})
	r.Get("/ping", h.Ping)
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetWholeMetric)
		r.Get("/{metricType}/{metricName}", h.GetMetric)
	})
	r.Get("/", h.AllMetrics)
	r.Post("/updates/", h.WriteAllMetrics)
	return r
}
