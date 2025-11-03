package access

import "net/http"

func init() {
	// wip
	http.HandleFunc("/admin/staff", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/admin/verify", func(w http.ResponseWriter, r *http.Request) {})
}
