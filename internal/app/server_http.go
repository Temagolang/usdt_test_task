package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus
)

func newHTTPServer(addr string) *http.Server {
func newHTTPServer(addr string, reg *prometheus.Registry) *http.Server {
	mux.HandleFunc("/healthz", handleHealthz)


.HandleFunc("/healthz
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
