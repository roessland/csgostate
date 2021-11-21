package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/roessland/csgostate/internal/logger"
	"net/http"
	"time"
)

func Serve(log logger.Logger) {
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9309",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("metrics server listening to to %s", srv.Addr)
	log.Infow("metrics server closed", "err", srv.ListenAndServe())
}
