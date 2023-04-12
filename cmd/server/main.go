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
	savingSettings := storage.NewSavingSettings(srvConfig.storeInterval, srvConfig.storeFile, srvConfig.dbConnString)
	metStorage.Restore(srvConfig.restore, savingSettings)
	go func() {
		for {
			<-time.After(srvConfig.storeInterval)
			if err := metStorage.Save(savingSettings.Database, savingSettings.StoreFile); err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	r := routers(metStorage, srvConfig.key, srvConfig.dbConnString)
	er := http.ListenAndServe(srvConfig.addr, r)
	if er != nil {
		fmt.Println(er.Error())
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
