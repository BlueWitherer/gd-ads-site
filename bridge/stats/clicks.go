package stats

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/stats/views", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Views endpoint hit")
		w.Write([]byte("42"))
	})
}
