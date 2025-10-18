package access

import (
	"net/http"

	"service/log"
)

func init() {
	http.HandleFunc("/admins", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Admins API service pinged")
		header := w.Header()
		header.Set("Content-Type", "text/plain")

		if code, err := Restrict(r.RemoteAddr); err != nil {
			http.Error(w, err.Error(), code)
		} else {
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Methods", "GET")
			header.Set("Access-Control-Allow-Headers", "Content-Type")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong!"))
		}
	})

	http.HandleFunc("/admins/get", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting admin status for user...")
		header := w.Header()
		header.Set("Content-Type", "text/plain")

		if code, err := Restrict(r.RemoteAddr); err != nil {
			http.Error(w, err.Error(), code)
		} else {
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Methods", "GET")
			header.Set("Access-Control-Allow-Headers", "Content-Type")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong!"))
		}
	})

	http.HandleFunc("/admins/all", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting all admins...")
		header := w.Header()
		header.Set("Content-Type", "text/plain")

		if code, err := Restrict(r.RemoteAddr); err != nil {
			http.Error(w, err.Error(), code)
		} else {
			header.Set("Access-Control-Allow-Methods", "GET")
			header.Set("Access-Control-Allow-Headers", "Content-Type")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong!"))
		}
	})
}
