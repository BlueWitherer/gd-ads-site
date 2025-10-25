package proxy

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"service/log"
)

type Level struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

func init() {
	http.HandleFunc("/proxy/level", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Error("Failed to parse form: %s", err.Error())
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// Try both parameter names for compatibility
			levelID := r.FormValue("levelID")
			if levelID == "" {
				levelID = r.FormValue("level-id")
			}

			if levelID == "" {
				http.Error(w, "levelID is required", http.StatusBadRequest)
				return
			}

			log.Info("Proxying request for level ID: %s", levelID)
			formData := url.Values{}
			formData.Set("levelID", levelID)
			formData.Set("secret", "Wmfd2893gb7")

			req, err := http.NewRequest("POST", "https://www.boomlings.com/database/downloadGJLevel22.php", strings.NewReader(formData.Encode()))
			if err != nil {
				log.Error("Failed to create request: %s", err.Error())
				http.Error(w, "Failed to create request", http.StatusInternalServerError)
				return
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("User-Agent", "")

			// Make the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Error("Failed to proxy request: %s", err.Error())
				http.Error(w, "Failed to fetch level data", http.StatusInternalServerError)
				return
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error("Failed to read response: %s", err.Error())
				http.Error(w, "Failed to read level data", http.StatusInternalServerError)
				return
			}

			header.Set("Content-Type", "text/plain")

			w.WriteHeader(resp.StatusCode)
			w.Write(body)

			log.Debug("Successfully proxied level request")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
