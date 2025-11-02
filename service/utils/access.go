package utils

type ArgonUser struct {
	Account int    `json:"account_id"` // Player account ID
	Token   string `json:"authtoken"`  // Authorization token
}

type ArgonValidation struct {
	Valid bool   `json:"valid"` // If user is valid
	Cause string `json:"cause"` // Cause for invalidation if any
}
