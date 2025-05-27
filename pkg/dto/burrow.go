package dto

import (
	"gophernet/pkg/db/ent"
)

// BurrowDto represents the data transfer object for burrows
type BurrowDto struct {
	Name       string  `json:"name"`
	Depth      float64 `json:"depth"`
	Width      float64 `json:"width"`
	IsOccupied bool    `json:"occupied"`
	Age        int     `json:"age"`
}

// ParseToModel converts BurrowDto to ent.Burrow
func (b *BurrowDto) ParseToModel() *ent.Burrow {
	return &ent.Burrow{
		Name:       b.Name,
		Depth:      b.Depth,
		Width:      b.Width,
		IsOccupied: b.IsOccupied,
		Age:        b.Age,
	}
}

// BurrowResponse represents a burrow in the system
type BurrowResponse struct {
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
