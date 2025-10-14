package api

import (
	"net/http"
	"strconv"

	"service/log"
)

func init() {
	http.HandleFunc("/api/ad", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting random ad...")
		header := w.Header()
		header.Set("Content-Type", "text/plain")

		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		var adFolder string

		query := r.URL.Query()
		adTypeStr := query.Get("type")

		typeNum, err := strconv.ParseInt(adTypeStr, 10, 64)
		if err != nil {
			log.Error("Failed to get ad type: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Map type to number
		switch typeNum {
		case 1:
			adFolder = "banner"

		case 2:
			adFolder = "square"

		case 3:
			adFolder = "skyscraper"

		default:
			http.Error(w, "Invalid ad type", http.StatusBadRequest)
			return
		}

		log.Debug("Getting random %s type ad...", adFolder)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	})
}
