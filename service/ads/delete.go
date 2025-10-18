package ads

import (
	"net/http"
	"strconv"

	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/ads/delete", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Attempting to delete ad(s)...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "DELETE")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		header.Set("Content-Type", "application/json")

		idStr := r.URL.Query().Get("ad_id")
		if idStr == "" {
			http.Error(w, "Missing ad ID parameter", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ad ID parameter", http.StatusBadRequest)
			return
		}

		err = database.DeleteAdvertisement(id)
		if err != nil {
			log.Error("Failed to delete advertisement: %s", err.Error())
			http.Error(w, "Failed to delete advertisement", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","message":"Advertisement deleted successfully"}`))
	})
}
