// сервис для хранения метрик ОС и получения их значений
package main

import (
	"log"
	"net/http"

	"goAdvancedTpl/internal/fabric/onstart"
	"goAdvancedTpl/internal/fabric/storage/dbstorage"
	"goAdvancedTpl/internal/fabric/storage/filestorage"
	"goAdvancedTpl/internal/server/config"
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/servermiddleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
)

func main() {

	onstart.WriteMessage(BuildVersion, BuildDate, BuildCommit)

	srvConfig := config.Config()
	var h *handlers.APIHandler
	if srvConfig.DBConnString != "" {
		metStorage := dbstorage.NewDBStorage(srvConfig.DBConnString, srvConfig.Restore)
		h = handlers.NewAPIHandler(metStorage, srvConfig.Key)

	} else {
		metStorage := filestorage.NewFileStorage(srvConfig.StoreInterval, srvConfig.StoreFile, srvConfig.Restore)
		h = handlers.NewAPIHandler(metStorage, srvConfig.Key)
	}

	r := routers(h, srvConfig.CryptoKey)
	err := http.ListenAndServe(srvConfig.Addr, r)
	if err != nil {
		log.Println(err.Error())
	}

}

func routers(h *handlers.APIHandler, cryptoKey string) *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(servermiddleware.GzipHandle)
	r.Use(servermiddleware.Decryption(cryptoKey))
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
