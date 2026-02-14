package models

// Error model
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

// Health model
type Health struct {
	Status string `json:"status"`
}
