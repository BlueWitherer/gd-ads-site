package proxy

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"bridge/access"
	"bridge/log"
)

type Level struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

func init() {
	http.HandleFunc("/proxy/level", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		if code, err := access.Restrict(r.RemoteAddr); err != nil {
			http.Error(w, err.Error(), code)
		} else {
			header.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			header.Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
			} else if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			} else {
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

				req, err := http.NewRequest("POST", "https://www.boomlings.com/database/downloadGJLevel22.php", strings.NewReader(formData.Encode()))
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

				header.Set("Content-Type", "text/plain")

				w.WriteHeader(resp.StatusCode)
				w.Write(body)

				log.Debug("Successfully proxied level request")
			}
		}
	})
}
