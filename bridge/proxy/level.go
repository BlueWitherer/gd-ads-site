package proxy

import (
	"bridge/log"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	http.HandleFunc("/api/proxy/level", handleLevelProxy)
}

func handleLevelProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Error("Failed to parse form: " + err.Error())
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	levelID := r.FormValue("levelID")
	if levelID == "" {
		http.Error(w, "levelID is required", http.StatusBadRequest)
		return
	}

	log.Info("Proxying request for level ID: " + levelID)
	formData := url.Values{}
	formData.Set("levelID", levelID)
	formData.Set("secret", "Wmfd2893gb7")

	req, err := http.NewRequest("POST", "https://www.boomlings.com/database/downloadGJLevel22.php",
		strings.NewReader(formData.Encode()))
	if err != nil {
		log.Error("Failed to create request: " + err.Error())
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to proxy request: " + err.Error())
		http.Error(w, "Failed to fetch level data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response: " + err.Error())
		http.Error(w, "Failed to read level data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	log.Debug("Successfully proxied level request")
}
