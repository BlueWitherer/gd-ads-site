package utils

import (
	"fmt"
	"time"
)

type AdType string // Dimensions of the ad image

const ( // Table to save stats
	AdTypeBanner     AdType = "banner"     // Horizontal ads
	AdTypeSquare     AdType = "square"     // Square ads
	AdTypeSkyscraper AdType = "skyscraper" // Vertical ads
)

type AdEvent string // Table to save stats to

const (
	AdEventView  AdEvent = "views"  // For views
	AdEventClick AdEvent = "clicks" // For clicks
)

type StatBy string // Row to filter through

const (
	StatByViews  StatBy = "total_views"  // Filter stats by ad
	StatByClicks StatBy = "total_clicks" // Filter stats by user
)

// Database row for advertisements listing
type Ad struct {
	AdID       int64     `json:"ad_id"`                 // Advertisement ID
	UserID     string    `json:"user_id"`               // Owner Discord user ID
	LevelID    int64     `json:"level_id"`              // Geometry Dash level ID
	Type       int       `json:"type"`                  // Type of advertisement
	ViewCount  int       `json:"view_count,omitempty"`  // Total registered views
	ClickCount int       `json:"click_count,omitempty"` // Total registered clicks
	ImageURL   string    `json:"image_url"`             // URL to the advertisement image
	Created    time.Time `json:"created_at"`            // First created
	Expiry     int64     `json:"expiry"`                // Unix time of expiration
	Pending    bool      `json:"pending"`               // Under review
}

func AdTypeFromInt(t int) (AdType, error) {
	switch t {
	case 1:
		return AdTypeBanner, nil
	case 2:
		return AdTypeSquare, nil
	case 3:
		return AdTypeSkyscraper, nil

	default:
		return "", fmt.Errorf("invalid ad type")
	}
}

func IntFromAdType(t AdType) (int, error) {
	switch t {
	case AdTypeBanner:
		return 1, nil
	case AdTypeSquare:
		return 2, nil
	case AdTypeSkyscraper:
		return 3, nil

	default:
		return 0, fmt.Errorf("invalid ad type")
	}
}
