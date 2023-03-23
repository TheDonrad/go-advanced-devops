package main

import (
	"flag"
	"fmt"
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/storage"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/go-chi/chi/v5"
)

type serverConfig struct {
	addr          string
	storeInterval time.Duration
	storeFile     string
	restore       bool
}

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

func srvConfig() serverConfig {
	srvConfig := serverConfig{
		addr:          "127.0.0.1:8080",
		storeInterval: 300,
		storeFile:     "/tmp/devops-metrics-db.json",
		restore:       true,
	}
	srvConfig.setConfigEnv()
	srvConfig.setConfigEnv()

	return srvConfig
}

func (srvConfig *serverConfig) setConfigFlags() {

	flag.StringVar(&srvConfig.addr, "a", srvConfig.addr, "host to listen on")

	flag.StringVar(&srvConfig.storeFile, "f", srvConfig.storeFile, "file to store metrics")

	flag.Func("i", "store interval", func(flagValue string) error {
		srvConfig.storeInterval, _ = time.ParseDuration(flagValue)
		return nil
	})
	flag.BoolVar(&srvConfig.restore, "r", srvConfig.restore, "restore")

	flag.Parse()

}

func (srvConfig *serverConfig) setConfigEnv() {
	var cfg struct {
		Addr          string `env:"ADDRESS"`
		StoreInterval string `env:"STORE_INTERVAL"`
		StoreFile     string `env:"STORE_FILE"`
		Restore       string `env:"RESTORE"`
	}

	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(strings.TrimSpace(cfg.Addr)) != 0 {
		srvConfig.addr = cfg.Addr
	}

	if len(strings.TrimSpace(cfg.StoreInterval)) != 0 {
		srvConfig.storeInterval, err = time.ParseDuration(cfg.StoreInterval)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if len(strings.TrimSpace(cfg.Restore)) != 0 {
		srvConfig.restore, err = strconv.ParseBool(cfg.Restore)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
