package utils

import "time"

type User struct {
	ID          string    `json:"id"`           // Discord user ID
	Username    string    `json:"username"`     // Discord username
	TotalViews  int       `json:"total_views"`  // Total registered views on all ads
	TotalClicks int       `json:"total_clicks"` // Total registered clicks on all ads
	IsAdmin     bool      `json:"is_admin"`     // Active administrator status
	IsStaff     bool      `json:"is_staff"`     // Active staff status
	Verified    bool      `json:"verified"`     // Trusted status
	Banned      bool      `json:"banned"`       // Banned status
	Created     time.Time `json:"created_at"`   // First created
	Updated     time.Time `json:"updated_at"`   // Last updated
}
