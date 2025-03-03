package models

import (
	"github.com/google/uuid"
)

// Guest represents a wedding guest
type Guest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	FamilySide  string    `json:"family_side"`
	Hongbao     float64   `json:"hongbao"`
	TotalGuests int       `json:"total_guests"`
	RSVPStatus  string    `json:"rsvp_status"`
	RSVPToken   string    `json:"rsvp_token"`
}

// NewGuest initializes a new Guest with a UUID
func NewGuest(name, email, familySide string, totalGuests int) *Guest {
	return &Guest{
		ID:          uuid.New(), // Generate a new UUID
		Name:        name,
		Email:       email,
		FamilySide:  familySide,
		Hongbao:     0, // Default value, can be updated later
		TotalGuests: totalGuests,
		RSVPStatus:  "Pending",
		RSVPToken:   uuid.New().String(), // Generate a unique RSVP token
	}
}
