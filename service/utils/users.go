package utils

import (
	"time"
)

type User struct {
	ID          string    `json:"id"`           // Discord user ID
	Username    string    `json:"username"`     // Discord username
	AvatarURL   string    `json:"avatar_url"`   // Discord user avatar URL
	TotalViews  uint64    `json:"total_views"`  // Total registered views on all ads
	TotalClicks uint64    `json:"total_clicks"` // Total registered clicks on all ads
	IsAdmin     bool      `json:"is_admin"`     // Active administrator status
	IsStaff     bool      `json:"is_staff"`     // Active staff status
	Verified    bool      `json:"verified"`     // Trusted status
	Banned      bool      `json:"banned"`       // Banned status
	BoostCount  uint      `json:"boost_count"`  // Available ad boosts
	Created     time.Time `json:"created_at"`   // First created
	Updated     time.Time `json:"updated_at"`   // Last updated
}

type Announcement struct {
	ID      uint      `json:"id"`         // Announcement ID
	User    User      `json:"user"`       // Announcement author
	Title   string    `json:"title"`      // Announcement title
	Content string    `json:"content"`    // Announcement content
	Created time.Time `json:"created_at"` // Created timestamp
}
