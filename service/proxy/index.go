package proxy

import (
	"net/http"

	"service/access"
	"service/log"
)

func init() {
	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Boomlings Proxy service pinged")
		header := w.Header()
		header.Set("Content-Type", "text/plain")

		if code, err := access.Restrict(r.RemoteAddr); err != nil {
			http.Error(w, err.Error(), code)
		} else {
			header.Set("Access-Control-Allow-Methods", "GET")
			header.Set("Access-Control-Allow-Headers", "Content-Type")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong!"))
		}
	})
}
