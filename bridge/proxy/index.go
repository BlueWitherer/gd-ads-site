package proxy

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/proxy", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Boomlings Proxy service pinged")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	})
}
