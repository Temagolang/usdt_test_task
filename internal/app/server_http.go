package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newHTTPServer(addr string, reg *prometheus.Registry) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handleHealthz)
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
