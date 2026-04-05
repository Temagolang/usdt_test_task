package app

import (
	"net/http"
)

func newHTTPServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handleHealthz)

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
