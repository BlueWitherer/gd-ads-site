package ads

import (
	"encoding/json"
	"net/http"
	"service/access"
	"service/database"
	"service/log"
	"service/utils"
	"strconv"
)

func init() {
	http.HandleFunc("/ads/reports", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			// require login
			uid, err := access.GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			u, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get user: %s", err.Error())
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			}

			if !u.IsAdmin && !u.IsStaff {
				log.Error("User of ID %s is not admin or staff", u.ID)
				http.Error(w, "User is not admin or staff", http.StatusUnauthorized)
				return
			}

			// Default behavior: get user's own ads
			rows, err := database.ListAllReports()
			if err != nil {
				log.Error("Failed to list reports: %s", err.Error())
				http.Error(w, "Failed to list reports", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(rows); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ads/reports/action", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			header.Set("Content-Type", "application/json")

			// require login
			uid, err := access.GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			u, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get user: %s", err.Error())
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			}

			if !u.IsAdmin && !u.IsStaff {
				log.Error("User of ID %s is not admin or staff", u.ID)
				http.Error(w, "User is not admin or staff", http.StatusUnauthorized)
				return
			}

			query := r.URL.Query()

			idStr := query.Get("id")
			if idStr == "" {
				http.Error(w, "Missing ad ID parameter", http.StatusBadRequest)
				return
			}

			actionStr := query.Get("action")
			if actionStr == "" {
				http.Error(w, "Missing action parameter", http.StatusBadRequest)
				return
			}

			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Error("Invalid ad ID parameter: %s", err.Error())
				http.Error(w, "Invalid ad ID parameter", http.StatusBadRequest)
				return
			}

			action, err := strconv.Atoi(idStr)
			if err != nil {
				log.Error("Invalid ad ID parameter: %s", err.Error())
				http.Error(w, "Invalid ad ID parameter", http.StatusBadRequest)
				return
			}

			report, err := database.GetReport(id)
			if err != nil {
				log.Error("Failed to get report: %s", err.Error())
				http.Error(w, "Failed to get report", http.StatusInternalServerError)
				return
			}

			if action == int(utils.ReportActionDelete) {
				ad, err := database.DeleteAdvertisement(report.Ad.AdID)
				if err != nil {
					log.Error("Failed to delete reported advertisement: %s", err.Error())
					http.Error(w, "Failed to delete reported advertisement", http.StatusInternalServerError)
					return
				}

				log.Info("Deleted reported advertisement of ID %d", ad.AdID)
			} else if action == int(utils.ReportActionBan) {
				if u.IsAdmin {
					user, err := database.BanUser(report.Ad.UserID)
					if err != nil {
						log.Error("Failed to ban owner of reported advertisement: %s", err.Error())
						http.Error(w, "Failed to ban owner of reported advertisement", http.StatusInternalServerError)
						return
					}

					log.Info("Banned owner of ID %s of reported advertisement", user.ID)
				} else {
					log.Error("Staff user of ID %s does not have permission to ban through reports", u.ID)
					http.Error(w, "Staff does not have permission to ban through reports", http.StatusUnauthorized)
					return
				}
			} else {
				log.Error("Invalid report action")
				http.Error(w, "Invalid report action", http.StatusBadRequest)
				return
			}

			err = database.FinishReport(report)
			if err != nil {
				log.Error("Failed to finalize report action: %s", err.Error())
				http.Error(w, "Failed to finalize report action", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Took action with report successfully"))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ads/reports/reject", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			header.Set("Content-Type", "application/json")

			// require login
			uid, err := access.GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			u, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get user: %s", err.Error())
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			}

			if !u.IsAdmin && !u.IsStaff {
				log.Error("User of ID %s is not admin or staff", u.ID)
				http.Error(w, "User is not admin or staff", http.StatusUnauthorized)
				return
			}

			query := r.URL.Query()

			idStr := query.Get("id")
			if idStr == "" {
				log.Error("Missing ad ID parameter")
				http.Error(w, "Missing ad ID parameter", http.StatusBadRequest)
				return
			}

			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Error("Invalid ad ID parameter: %s", err.Error())
				http.Error(w, "Invalid ad ID parameter", http.StatusBadRequest)
				return
			}

			report, err := database.GetReport(id)
			if err != nil {
				log.Error("Failed to get report: %s", err.Error())
				http.Error(w, "Failed to get report", http.StatusInternalServerError)
				return
			}

			err = database.FinishReport(report)
			if err != nil {
				log.Error("Failed to finish report: %s", err.Error())
				http.Error(w, "Failed to finish report", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Rejected report successfully"))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
