package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/g4l1l10/rsvp-backend/models"
	"github.com/g4l1l10/rsvp-backend/repository"
	"github.com/g4l1l10/rsvp-backend/utils"

	"github.com/google/uuid"
)

// GuestService defines business logic for guest management
type GuestService struct {
	Repo *repository.GuestRepository
}

// NewGuestService initializes a new guest service
func NewGuestService(repo *repository.GuestRepository) *GuestService {
	return &GuestService{Repo: repo}
}

// AddGuest validates input and creates a new guest
func (s *GuestService) AddGuest(name, email, familySide string, totalGuests int) (*models.Guest, error) {
	// Validate inputs
	if name == "" || email == "" || familySide == "" || totalGuests <= 0 {
		return nil, errors.New("invalid input: all fields must be provided and total guests must be greater than zero")
	}

	// Create a new guest with UUID and unique RSVP token
	guest := models.NewGuest(name, email, familySide, totalGuests)

	// Store guest in database
	err := s.Repo.CreateGuest(guest)
	if err != nil {
		return nil, fmt.Errorf("failed to add guest: %v", err)
	}

	log.Println("📧 Sending email invitation to:", email)

	// Send Email (Run in a goroutine to prevent blocking)
	go func() {
		err = utils.SendEmail(guest.Name, guest.Email, guest.RSVPToken)
		if err != nil {
			log.Println("⚠️ Warning: Guest added but email not sent:", err)
		}
	}()

	return guest, nil
}

// GetAllGuests retrieves all guests
func (s *GuestService) GetAllGuests() ([]models.Guest, error) {
	guests, err := s.Repo.GetAllGuests()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch guests: %v", err)
	}
	return guests, nil
}

// GetGuestByID retrieves a guest by UUID
func (s *GuestService) GetGuestByID(id uuid.UUID) (*models.Guest, error) {
	guest, err := s.Repo.GetGuestByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving guest: %v", err)
	}
	return guest, nil
}

// GetGuestByToken retrieves a guest by their RSVP token
func (s *GuestService) GetGuestByToken(token string) (*models.Guest, error) {
	guest, err := s.Repo.GetGuestByToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid RSVP token: %v", err)
	}
	return guest, nil
}

// GetGuestByEmail retrieves a guest by email
func (s *GuestService) GetGuestByEmail(email string) (*models.Guest, error) {
	guest, err := s.Repo.GetGuestByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error retrieving guest by email: %v", err)
	}
	return guest, nil
}

// GetGuestsByRSVP retrieves all guests with a specific RSVP status
func (s *GuestService) GetGuestsByRSVP(status string) ([]models.Guest, error) {
	guests, err := s.Repo.GetGuestByRSVP(status)
	if err != nil {
		return nil, fmt.Errorf("error retrieving guests with RSVP status %s: %v", status, err)
	}
	return guests, nil
}

// UpdateGuest updates an existing guest's details
func (s *GuestService) UpdateGuest(guest *models.Guest) error {
	// Validate guest data before updating
	if guest.ID == uuid.Nil {
		return errors.New("invalid guest ID")
	}

	err := s.Repo.UpdateGuest(guest)
	if err != nil {
		return fmt.Errorf("failed to update guest: %v", err)
	}
	return nil
}

// DeleteGuest removes a guest from the system
func (s *GuestService) DeleteGuest(id uuid.UUID) error {
	err := s.Repo.DeleteGuest(id)
	if err != nil {
		return fmt.Errorf("failed to delete guest: %v", err)
	}
	return nil
}
