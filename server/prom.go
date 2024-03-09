package server

import (
	"github.com/504dev/logr/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func ListenPROM() error {
	address := config.Get().Bind.Prom
	if address == "" {
		return nil
	}
	//use separated ServeMux to prevent handling on the global Mux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(address, mux)
}
