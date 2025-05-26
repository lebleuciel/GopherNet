package models

// Burrow represents a burrow in the system
type Burrow struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Depth      float64 `json:"depth"`
	Width      float64 `json:"width"`
	IsOccupied bool    `json:"is_occupied"`
	Age        int     `json:"age"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error string `json:"error"`
}
