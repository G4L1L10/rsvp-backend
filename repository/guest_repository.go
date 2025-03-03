package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/g4l1l10/rsvp-backend/models"

	"github.com/google/uuid"
)

// GuestRepository handles database operations for guests
type GuestRepository struct {
	DB *sql.DB
}

// NewGuestRepository initializes a new repository instance
func NewGuestRepository(db *sql.DB) *GuestRepository {
	return &GuestRepository{DB: db}
}

// CreateGuest inserts a new guest into the database securely
func (r *GuestRepository) CreateGuest(guest *models.Guest) error {
	query := `
		INSERT INTO guests (id, name, email, family_side, hongbao, total_guests, rsvp_status, rsvp_token)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`
	err := r.DB.QueryRow(query, guest.ID, guest.Name, guest.Email, guest.FamilySide, guest.Hongbao, guest.TotalGuests, guest.RSVPStatus, guest.RSVPToken).Scan(&guest.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllGuests retrieves all guests from the database
func (r *GuestRepository) GetAllGuests() ([]models.Guest, error) {
	query := "SELECT id, name, email, family_side, hongbao, total_guests, rsvp_status, rsvp_token FROM guests"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guests []models.Guest
	for rows.Next() {
		var g models.Guest
		if err := rows.Scan(&g.ID, &g.Name, &g.Email, &g.FamilySide, &g.Hongbao, &g.TotalGuests, &g.RSVPStatus, &g.RSVPToken); err != nil {
			return nil, err
		}
		guests = append(guests, g)
	}

	return guests, nil
}

// GetGuestByID fetches a single guest securely using a UUID
func (r *GuestRepository) GetGuestByID(id uuid.UUID) (*models.Guest, error) {
	query := "SELECT id, name, email, family_side, hongbao, total_guests, rsvp_status, rsvp_token FROM guests WHERE id = $1"
	row := r.DB.QueryRow(query, id)

	var guest models.Guest
	err := row.Scan(&guest.ID, &guest.Name, &guest.Email, &guest.FamilySide, &guest.Hongbao, &guest.TotalGuests, &guest.RSVPStatus, &guest.RSVPToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("guest not found")
		}
		return nil, err
	}

	return &guest, nil
}

// GetGuestByToken fetches a guest using their unique RSVP token
func (r *GuestRepository) GetGuestByToken(token string) (*models.Guest, error) {
	query := "SELECT id, name, email, family_side, hongbao, total_guests, rsvp_status, rsvp_token FROM guests WHERE rsvp_token = $1"
	row := r.DB.QueryRow(query, token)

	var guest models.Guest
	err := row.Scan(&guest.ID, &guest.Name, &guest.Email, &guest.FamilySide, &guest.Hongbao, &guest.TotalGuests, &guest.RSVPStatus, &guest.RSVPToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid RSVP token")
		}
		return nil, err
	}

	return &guest, nil
}

// GetGuestByEmail fetches a guest using their email
func (r *GuestRepository) GetGuestByEmail(email string) (*models.Guest, error) {
	query := "SELECT id, name, email, family_side, hongbao, total_guests, rsvp_status, rsvp_token FROM guests WHERE email = $1"
	row := r.DB.QueryRow(query, email)

	var guest models.Guest
	err := row.Scan(&guest.ID, &guest.Name, &guest.Email, &guest.FamilySide, &guest.Hongbao, &guest.TotalGuests, &guest.RSVPStatus, &guest.RSVPToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("guest not found with provided email")
		}
		return nil, err
	}

	return &guest, nil
}

// GetGuestByRSVP retrieves guests based on their RSVP status (e.g., "Attending", "Not Attending", "Pending")
func (r *GuestRepository) GetGuestByRSVP(status string) ([]models.Guest, error) {
	query := "SELECT id, name, email, family_side, hongbao, total_guests, rsvp_status, rsvp_token FROM guests WHERE rsvp_status = $1"
	rows, err := r.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guests []models.Guest
	for rows.Next() {
		var g models.Guest
		if err := rows.Scan(&g.ID, &g.Name, &g.Email, &g.FamilySide, &g.Hongbao, &g.TotalGuests, &g.RSVPStatus, &g.RSVPToken); err != nil {
			return nil, err
		}
		guests = append(guests, g)
	}

	return guests, nil
}

// UpdateGuest updates a guest's information securely
func (r *GuestRepository) UpdateGuest(guest *models.Guest) error {
	// Fetch the existing guest details
	var existingGuest models.Guest
	query := "SELECT name, email, family_side, hongbao, total_guests, rsvp_status, rsvp_token FROM guests WHERE id = $1"
	err := r.DB.QueryRow(query, guest.ID).Scan(
		&existingGuest.Name,
		&existingGuest.Email,
		&existingGuest.FamilySide,
		&existingGuest.Hongbao,
		&existingGuest.TotalGuests,
		&existingGuest.RSVPStatus,
		&existingGuest.RSVPToken,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch existing guest: %v", err)
	}

	// Preserve existing values if fields are not provided in the request
	if guest.Name == "" {
		guest.Name = existingGuest.Name
	}
	if guest.Email == "" {
		guest.Email = existingGuest.Email
	}
	if guest.FamilySide == "" {
		guest.FamilySide = existingGuest.FamilySide
	}
	if guest.RSVPStatus == "" {
		guest.RSVPStatus = existingGuest.RSVPStatus
	}
	if guest.RSVPToken == "" {
		guest.RSVPToken = existingGuest.RSVPToken
	}

	// Update query
	updateQuery := `
		UPDATE guests
		SET name = $1, email = $2, family_side = $3, hongbao = $4, total_guests = $5, rsvp_status = $6, rsvp_token = $7
		WHERE id = $8;
	`
	_, err = r.DB.Exec(updateQuery, guest.Name, guest.Email, guest.FamilySide, guest.Hongbao, guest.TotalGuests, guest.RSVPStatus, guest.RSVPToken, guest.ID)
	if err != nil {
		return fmt.Errorf("failed to update guest: %v", err)
	}

	return nil
}

// DeleteGuest removes a guest securely from the database
func (r *GuestRepository) DeleteGuest(id uuid.UUID) error {
	query := "DELETE FROM guests WHERE id = $1"
	_, err := r.DB.Exec(query, id)
	return err
}
