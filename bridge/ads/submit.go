package ads

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/ads/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Parse form with 10MB limit
			r.ParseMultipartForm(10 << 20)

			// Get image file
			file, handler, err := r.FormFile("image-upload")
			if err != nil {
				http.Error(w, "Image not found", http.StatusBadRequest)
				return
			}

			defer file.Close()

			var adFolder string = r.Form.Get("type")

			// Create target folder
			targetDir := filepath.Join("..", "..", "ads", adFolder)
			os.MkdirAll(targetDir, os.ModePerm)

			// Save file
			dstPath := filepath.Join(targetDir, handler.Filename)
			dst, err := os.Create(dstPath)
			if err != nil {
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}

			defer dst.Close()

			io.Copy(dst, file)

			log.Info("Saved image to " + dstPath)
			w.Write([]byte(`{"status":"image saved"}`))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
