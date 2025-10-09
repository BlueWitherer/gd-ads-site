package ads

import (
	"net/http"

	"service/access"
	"service/log"
)

func init() {
	http.HandleFunc("/ads/delete", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Attempting to delete ad(s)...")
		log.Warn("This feature has not been implemented yet!")
		header := w.Header()
		header.Set("Content-Type", "application/json")

		if code, err := access.Restrict(r.RemoteAddr); err != nil {
			http.Error(w, err.Error(), code)
		} else {
			header.Set("Access-Control-Allow-Methods", "POST")
			header.Set("Access-Control-Allow-Headers", "Content-Type")

			w.WriteHeader(http.StatusNotImplemented)
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	})
}
