package utils

import "time"

type ArgonUser struct {
	Account      int       `json:"account_id"`    // Player account ID
	Token        string    `json:"authtoken"`     // Authorization token
	ReportBanned bool      `json:"report_banned"` // Is banned from reporting
	ValidAt      time.Time `json:"valid_at"`      // Token last validated at
}

type ArgonValidation struct {
	Valid bool   `json:"valid"` // If user is valid
	Cause string `json:"cause"` // Cause for invalidation if any
}
