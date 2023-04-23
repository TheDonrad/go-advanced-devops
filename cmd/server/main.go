package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"goAdvancedTpl/internal/fabric/storage/dbstorage"
	"goAdvancedTpl/internal/fabric/storage/filestorage"
	"goAdvancedTpl/internal/server/config"
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/servermiddleware"
	"log"
	"net/http"
)

func main() {

	srvConfig := config.SrvConfig()
	var h *handlers.APIHandler
	if srvConfig.DBConnString != "" {
		metStorage := dbstorage.NewDBStorage(srvConfig.DBConnString, srvConfig.Restore)
		h = handlers.NewAPIHandler(metStorage, srvConfig.Key)

	} else {
		metStorage := filestorage.NewFileStorage(srvConfig.StoreInterval, srvConfig.StoreFile, srvConfig.Restore)
		h = handlers.NewAPIHandler(metStorage, srvConfig.Key)
	}

	r := routers(h)
	err := http.ListenAndServe(srvConfig.Addr, r)
	if err != nil {
		log.Println(err.Error())
	}

}

func routers(h *handlers.APIHandler) *chi.Mux {

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
