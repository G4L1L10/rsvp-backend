package handlers

import (
	"log"
	"net/http"

	"github.com/g4l1l10/rsvp-backend/service"

	"github.com/g4l1l10/rsvp-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GuestHandler handles business logic for guest operations
type GuestHandler struct {
	Service *service.GuestService
}

// NewGuestHandler initializes a new guest handler
func NewGuestHandler(service *service.GuestService) *GuestHandler {
	return &GuestHandler{Service: service}
}

// AddGuest handles adding a new guest
func (h *GuestHandler) AddGuest(ctx *gin.Context) {
	log.Println("📥 Received request to add guest")

	var req struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		FamilySide  string `json:"family_side" binding:"required"`
		TotalGuests int    `json:"total_guests" binding:"required,gt=0"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("❌ Error binding request:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	guest, err := h.Service.AddGuest(req.Name, req.Email, req.FamilySide, req.TotalGuests)
	if err != nil {
		log.Println("❌ Failed to add guest:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("✅ Guest successfully added:", guest.ID)
	ctx.JSON(http.StatusCreated, guest)
}

// GetAllGuests retrieves all guests
func (h *GuestHandler) GetAllGuests(ctx *gin.Context) {
	guests, err := h.Service.GetAllGuests()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, guests)
}

// GetGuestByID retrieves a guest by their UUID
func (h *GuestHandler) GetGuestByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid guest ID"})
		return
	}

	guest, err := h.Service.GetGuestByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, guest)
}

// GetGuestByEmail retrieves a guest by their email
func (h *GuestHandler) GetGuestByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	guest, err := h.Service.GetGuestByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, guest)
}

// GetGuestByToken retrieves a guest by their RSVP token
func (h *GuestHandler) GetGuestByToken(ctx *gin.Context) {
	token := ctx.Param("token")
	guest, err := h.Service.GetGuestByToken(token)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, guest)
}

// UpdateGuest updates a guest's information
func (h *GuestHandler) UpdateGuest(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid guest ID"})
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Email       string  `json:"email"`
		FamilySide  string  `json:"family_side"`
		Hongbao     float64 `json:"hongbao"`
		TotalGuests int     `json:"total_guests"`
		RSVPStatus  string  `json:"rsvp_status"`
	}

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	guest := &models.Guest{
		ID:          id,
		Name:        req.Name,
		Email:       req.Email,
		FamilySide:  req.FamilySide,
		Hongbao:     req.Hongbao,
		TotalGuests: req.TotalGuests,
		RSVPStatus:  req.RSVPStatus,
	}

	err = h.Service.UpdateGuest(guest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "guest updated successfully"})
}

// DeleteGuest removes a guest
func (h *GuestHandler) DeleteGuest(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid guest ID"})
		return
	}

	err = h.Service.DeleteGuest(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "guest deleted successfully"})
}
