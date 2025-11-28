package proxy

import (
	"fmt"
	"net/http"

	"service/log"
)

func init() {
	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Boomlings Proxy service pinged")
		header := w.Header()
		header.Set("Content-Type", "text/plain")

		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong!")
	})
}
