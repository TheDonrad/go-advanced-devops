package main

import (
	"log"
	"net/http"

	"goAdvancedTpl/internal/fabric/storage/dbstorage"
	"goAdvancedTpl/internal/fabric/storage/filestorage"
	"goAdvancedTpl/internal/server/config"
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/servermiddleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	srvConfig := config.SrvConfig()
	var metStorage handlers.Storage
	if srvConfig.DBConnString != "" {
		metStorage = dbstorage.NewDBStorage(srvConfig.DBConnString)

	} else {
		metStorage = filestorage.NewFileStorage(srvConfig.StoreInterval, srvConfig.StoreFile)
	}

	if srvConfig.Restore {
		metStorage.Restore()
	}

	r := routers(metStorage, srvConfig.Key)
	err := http.ListenAndServe(srvConfig.Addr, r)
	if err != nil {
		log.Println(err.Error())
	}

}

func routers(metStorage handlers.Storage, key string) *chi.Mux {

	h := handlers.NewAPIHandler(metStorage, key)
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
